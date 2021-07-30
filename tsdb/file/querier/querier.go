package querier

import (
	"fmt"
	"sort"
	"time"

	"github.com/bench-routes/bench-routes/src/lib/utils/decode"
	tsdb "github.com/bench-routes/bench-routes/tsdb/file"
)

type Query struct {
	dbPath    string
	typ       uint8
	rang      queryRange
	timestamp time.Time
}

// queryRange is the range of time within which the set of blocks
// are to be looked for.
type queryRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// QueryResponse is the response sent after processing the query.
type QueryResponse struct {
	TimeSeriesPath string         `json:"timeSeriesPath"`
	Type           string         `json:"type"`
	Range          queryRange     `json:"range"`
	EvaluationTime string         `json:"evaluationTime"`
	Values         []DecodedBlock `json:"values"`
}

type DecodedBlock struct {
	Value          interface{} `json:"value"`
	Timestamp      string      `json:"timestamp"`
	NormalizedTime int64       `json:"normalizedTime"`
}

const (
	// TypeRange corresponds to the querier that requires querying over a
	// time range.
	TypeRange uint8 = 0
	// TypeFirst is a querier that requires the first sample only.
	TypeFirst uint8 = 1
)

// New Creates a new Query which can be executed.
func New(typ uint8, path string, start int64, end int64) (*Query, error) {
	rang := queryRange{
		Start: start,
		End:   end,
	}
	query := &Query{
		dbPath: path,
		rang:   rang,
		typ:    typ,
	}
	err := query.Validate()
	if err != nil {
		return nil, fmt.Errorf("validation error: %s", err)
	}
	return query, nil
}

//Exec executes the query for the given range and returns the QueryResponse.
func (q *Query) Exec() (*QueryResponse, error) {
	bstream, err := tsdb.FetchChainStream(q.dbPath)
	if err != nil {
		return nil, fmt.Errorf("loading chain error: %s", err)
	}
	q.timestamp = time.Now()
	sortedStream := sortStream(bstream)
	resStream, err := q.fetchBlocks(sortedStream)
	if err != nil {
		return nil, fmt.Errorf("query exec error : %s", err)
	}
	if len(resStream) == 0 {
		return &QueryResponse{
			TimeSeriesPath: q.dbPath,
			Values:         []DecodedBlock{},
			Type:           "",
			Range:          q.rang,
			EvaluationTime: time.Since(q.timestamp).String(),
		}, nil
	}
	typ := fetchType(resStream)

	// decoding blockstream according to block type
	blocksDecoder := decode.NewBlockDecoding(typ)
	var decodedStream []DecodedBlock
	for _, block := range resStream {
		decodedStream = append(decodedStream, DecodedBlock{
			Value:          blocksDecoder.Decode(block),
			Timestamp:      block.Timestamp,
			NormalizedTime: block.NormalizedTime,
		})
	}

	return &QueryResponse{
		TimeSeriesPath: q.dbPath,
		Type:           blocksDecoder.Type,
		Values:         decodedStream,
		Range:          q.rang,
		EvaluationTime: time.Since(q.timestamp).String(),
	}, nil
}

//fetchBlocks fetches the blocks from the stream within the query range using binary search.
func (q *Query) fetchBlocks(stream *[]tsdb.Block) ([]tsdb.Block, error) {
	streamLen := len(*stream)
	if streamLen == 0 {
		return nil, fmt.Errorf("no blocks found in chain path")
	}

	// Here firstBlockIndex represents the oldest block index(according to time) for the
	// given query range and vice-versa for lastBlockIndex.
	// Also firstBlockTime represents the oldest block time in the stream
	// and lastBlockTime represents the latest block time in the stream.
	var firstBlockIndex, lastBlockIndex int
	firstBlockTime, lastBlockTime := (*stream)[0].NormalizedTime, (*stream)[streamLen-1].NormalizedTime
	switch q.typ {
	case TypeFirst:
		return []tsdb.Block{(*stream)[streamLen-1]}, nil
	case TypeRange:
		if q.rang.Start < firstBlockTime {
			firstBlockIndex = 0
		} else if q.rang.Start > lastBlockTime {
			return []tsdb.Block{}, nil
		} else {
			firstBlockIndex = binSearch(stream, q.rang.Start, 0, streamLen-1)
		}

		if q.rang.End > lastBlockTime {
			lastBlockIndex = streamLen - 1
		} else if q.rang.End < firstBlockTime {
			return []tsdb.Block{}, nil
		} else {
			lastBlockIndex = binSearch(stream, q.rang.End, 0, streamLen-1)
		}
	default:
		return nil, fmt.Errorf("typ error: invalid query type")
	}
	return (*stream)[firstBlockIndex : lastBlockIndex+1], nil
}

// Validate validates the query.
func (q *Query) Validate() error {
	if q.rang.Start > q.rang.End {
		return fmt.Errorf("range error: start time is greater than end time")
	}
	if ok := tsdb.VerifyChainPathExists(q.dbPath); !ok {
		return fmt.Errorf("dbpath error: path doesn't exists")
	}

	if q.typ != TypeFirst && q.typ != TypeRange {
		return fmt.Errorf("typ error: invalid query type")
	}
	return nil
}

// binSearch searches for the index of the block in the stream in O(log(n))
func binSearch(stream *[]tsdb.Block, time int64, startIndex int, endIndex int) int {
	if startIndex >= endIndex {
		return startIndex
	}

	mid := (startIndex + endIndex) / 2
	if (*stream)[mid].NormalizedTime == time {
		return mid
	} else if (*stream)[mid].NormalizedTime > time {
		return binSearch(stream, time, startIndex, mid-1)
	} else {
		return binSearch(stream, time, mid+1, endIndex)
	}
}

type blockStream []tsdb.Block

func (a blockStream) Len() int {
	return len(a)
}
func (a blockStream) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a blockStream) Less(i, j int) bool {
	return a[i].NormalizedTime < a[j].NormalizedTime
}

func sortStream(stream []tsdb.Block) *[]tsdb.Block {
	sort.Sort(blockStream(stream))
	return &stream
}

func fetchType(stream []tsdb.Block) string {
	for _, b := range stream {
		if b.Type != "null" {
			return b.Type
		}
	}
	return "null"
}
