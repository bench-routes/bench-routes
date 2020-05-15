package querier

import (
	"math"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/utils/decode"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

const (
	// TypeRange corresponds to the querier that requires querying over a
	// time range.
	TypeRange uint8 = 0
	// TypeFirst is a querier that requires the first sample only.
	TypeFirst uint8 = 1
	// TypeLast is a querier that requires the last sample only.
	TypeLast uint8 = 2
)

// Querier is a querying type.
type Querier struct {
	DBPath        string
	CollectorPath string
	Type          uint8
}

// Query is a complex of db path, type and range of time-series.
type Query struct {
	Path          string
	Range         *queryRange
	stamp         time.Time
	encoding      bool
	queryResponse QueryResponse
	Type          uint8
}

// queryRange is the range of time within which the set of blocks
// are to be looked for.
type queryRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// New returns a new querier that implements querying.
func New(TSDBPath, collectorPath string, Type uint8) *Querier {
	return &Querier{
		DBPath:        TSDBPath,
		CollectorPath: collectorPath,
		Type:          Type,
	}
}

// QueryBuilder builds a query that can be executed over a time-series
// range.
// The path of the presence of the time-series values will be formed as:
// <DBPath>/<ofType>/chunk_<middle>/<url>.json
func (qs *Querier) QueryBuilder() *Query {
	return &Query{
		Path: qs.DBPath,
		// range "nil" represents all time-series as querying
		Range: nil,
		Type:  qs.Type,
	}
}

// SetRange is used to set the range of the query.
// Note that fromTimestamp > tillTimestamp denoting the range
// [tillTimestamp, fromTimestamp] in the axis of time, [] denoting
// the including value.
// If no range is set, the query is assumed to be covering all the
// time-series present in that file (or db containing that time-series).
func (q *Query) SetRange(fromTimestamp, tillTimestamp int64) {
	q.Range = &queryRange{
		Start: fromTimestamp,
		End:   tillTimestamp,
	}
}

// Exec executes the Query returned from the QueryBuilder.
func (q *Query) Exec() []byte {
	chainReadOnly := tsdb.ReadOnly(q.Path).Refresh()
	bstream := chainReadOnly.BlockStream()

	data, ok := q.exec(*bstream, true).([]byte)
	if !ok {
		logger.Terminal("p", "invalid []byte extracting from interface{}")
	}
	return data
}

// ExecWithoutEncode executes the Query without encoding the result to []byte,
// rather keeps the result as default QueryResponse.
func (q *Query) ExecWithoutEncode() QueryResponse {
	chainReadOnly := tsdb.ReadOnly(q.Path).Refresh()
	bstream := chainReadOnly.BlockStream()
	data, ok := q.exec(*bstream, false).(QueryResponse)
	if !ok {
		logger.Terminal("p", "invalid []byte extracting from interface{}")
	}
	return data
}

func (q *Query) exec(blockStream []tsdb.Block, encoding bool) interface{} {
	// start represents the starting of the timer that will calculate
	// the total time involved in performing a particular query.
	// This can be later benchmarked and compared with other algorithms.
	base, stamp := getBaseResponse(q.Range)
	q.encoding = encoding
	q.stamp = stamp
	q.queryResponse = base
	if len(blockStream) == 0 {
		return q.ReturnMessageResponse("EMPTY_SAMPLES")
	}

	var (
		lengthBlockStream = len(blockStream)
		// since we are querying the time-series, the first block is the most
		// recently added block into the time-series. This obviously
		// corresponds to the last block in the block stream. Similarly, the last time will
		// correspond to the oldest block which is the first block in the
		// block stream.
		timeFirstBlock      = blockStream[lengthBlockStream-1].GetNormalizedTime()
		timeLastBlock       = blockStream[0].GetNormalizedTime()
		onEndTimestamp      = false
		decodedBlockStream  []interface{}
		resultingBlockSlice []tsdb.Block
	)
	issues := q.validate(timeFirstBlock, timeLastBlock, lengthBlockStream)
	if issues != nil {
		return issues
	}

	switch q.Type {
	case TypeFirst:
		resultingBlockSlice = []tsdb.Block{blockStream[lengthBlockStream-1]}
	case TypeLast:
		resultingBlockSlice = []tsdb.Block{blockStream[0]}
	case TypeRange:
		// A nil range represents to return all time-series value as response.
		if q.Range == nil {
			base = QueryResponse{
				TimeInvolved: time.Since(stamp),
				Range: queryRange{
					Start: int64(math.MaxInt64),
					End:   int64(math.MinInt64),
				},
				Value: blockStream,
			}
			return encode(base, encoding)
		}
		if len(blockStream) == 0 {
			return q.ReturnMessageResponse("NO_BLOCKS_IN_MENTIONED_DB_PATH")
		}

		var (
			startPos, endPos = lengthBlockStream - 1, 0
			startPosFound    bool
		)

		for i := lengthBlockStream - 1; i >= 0; i-- {
			block := blockStream[i]
			if block.GetNormalizedTime() <= q.Range.Start && !startPosFound {
				startPos = i
				startPosFound = true
				continue
			}
			// When the input end time has a time value much less than the
			// normalized-time of the last block in the time-series, we
			// need to prevent slicing the last most bock caused by
			// endPos + 1.
			if block.GetNormalizedTime() == timeLastBlock || block.GetNormalizedTime() == q.Range.End {
				onEndTimestamp = true
			}
			if block.GetNormalizedTime() <= q.Range.End && startPosFound {
				endPos = i
				break
			}
		}
		// we need to skip the first block in the range in order to fall inside
		// the range of [Start, End] query range. If not, this will lead to
		// Start - 1, End range of blocks.
		resultingBlockSlice = blockStream[endPos+1 : startPos+1]
		if onEndTimestamp {
			resultingBlockSlice = blockStream[endPos : startPos+1]
		}
	}

	if len(resultingBlockSlice) == 0 {
		return q.ReturnNILResponse()
	}
	// decode the selected range of blocks
	blocksDecoder := decode.NewBlockDecoding(resultingBlockSlice[0].Type)
	for i := range resultingBlockSlice {
		decodedBlockStream = append(decodedBlockStream, queryValue{
			Timestamp:      resultingBlockSlice[i].Timestamp,
			Value:          blocksDecoder.Decode(resultingBlockSlice[i]),
			NormalizedTime: resultingBlockSlice[i].GetNormalizedTime(),
		})
	}
	base = QueryResponse{
		TimeSeriesPath: q.Path,
		TimeInvolved:   time.Since(stamp),
		Range:          *q.Range,
		Value:          decodedBlockStream,
	}
	return encode(base, encoding)
}

// validate performs pre-validations on the query in order to avoid faults
// while traversing through the time-series values. It also returns the default
// responses according to the case.
func (q *Query) validate(timeFirstBlock, timeLastBlock int64, l int) []byte {
	if q.Range == nil {
		return nil
	}
	if q.Range.Start < q.Range.End {
		return q.ReturnMessageResponse("ERROR_fromTimestamp_LESS_THAN_tillTimestamp")
	}
	if q.Range.Start < timeLastBlock || q.Range.End > timeFirstBlock {
		return q.ReturnNILResponse()
	}
	if l == 0 {
		return q.ReturnNILResponse()
	}
	return nil
}

// ReturnNILResponse is used only in cases when the response is known to be
// null value.
func (q *Query) ReturnNILResponse() []byte {
	q.queryResponse.TimeInvolved = time.Since(q.stamp)
	data, ok := encode(q.queryResponse, q.encoding).([]byte)
	if !ok {
		logger.Terminal("p", "invalid []byte extracting from interface{}")
	}
	return data
}

// ReturnMessageResponse is used only for sending simple messages. This should
// not be used to send JSON based results.
func (q *Query) ReturnMessageResponse(message string) []byte {
	q.queryResponse = QueryResponse{
		TimeInvolved: time.Since(q.stamp),
		Value:        message,
	}
	data, ok := encode(q.queryResponse, q.encoding).([]byte)
	if !ok {
		logger.Terminal("p", "invalid []byte extracting from interface{}")
	}
	return data
}
