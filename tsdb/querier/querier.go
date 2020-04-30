package querier

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

// Querier is a querying type.
type Querier struct {
	writer        http.ResponseWriter
	DBPath        string
	CollectorPath string
}

// Query is a complex of db path, type and range of time-series.
type Query struct {
	ws            http.ResponseWriter
	Path          string
	Range         *queryRange
	stamp         time.Time
	queryResponse response
}

// queryRange is the range of time within which the set of blocks
// are to be looked for.
type queryRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// New returns a new Querier that implements querying.
func New(responseWriter http.ResponseWriter, TSDBPath, collectorPath string) *Querier {
	return &Querier{
		writer:        responseWriter,
		DBPath:        TSDBPath,
		CollectorPath: collectorPath,
	}
}

// QueryBuilder builds a query that can be executed over a time-series
// range.
// The path of the presence of the time-series values will be formed as:
// <DBPath>/<ofType>/chunk_<middle>/<url>.json
func (qs *Querier) QueryBuilder() *Query {
	return &Query{
		ws:   qs.writer,
		Path: qs.DBPath,
		// range "nil" represents all time-series as querying
		Range: nil,
	}
}

// SetRange is used to set the range of the query.
// Note that fromTimestamp > tillTimestamp denotingthe range
// [tillTimestamp, fromTimestamp] in the axis of time, [] denoting
// the invlucing value.
// If no range is set, the query is assumed to be covering all the
// time-series present in that file (or db containing that time-series).
func (q *Query) SetRange(fromTimestamp, tillTimestamp int64) {
	q.Range = &queryRange{
		Start: fromTimestamp,
		End:   tillTimestamp,
	}
}

// Exec executes the Query returned from the QueryBuilder
func (q *Query) Exec() {
	chainReadOnly := tsdb.ReadOnly(q.Path).Refresh()
	bstream := chainReadOnly.BlockStream()

	q.exec(*bstream)
}

func (q *Query) exec(blockStream []tsdb.Block) {
	var js []byte

	// start represents the starting of the timer that will calculate
	// the total time involved in performing a particular query.
	// This can be later benchmarked and compared with other algorithms.
	base, stamp := getBaseResponse(*q.Range)
	q.stamp = stamp
	q.queryResponse = base

	var (
		lengthBlockStream = len(blockStream)
		// since we are querying the time-series, the first block is the most
		// recently added block into the time-series. This obviously
		// corresponds to the last block in the blockstream. Similarly, the last time will
		// correspond to the oldest block which is the first block in the
		// blockstream.
		timeFirstBlock = blockStream[lengthBlockStream-1].GetNormalizedTime()
		timeLastBlock  = blockStream[0].GetNormalizedTime()
	)
	status := q.validate(timeFirstBlock, timeLastBlock)
	if !status {
		logger.Terminal(fmt.Errorf("invalid query received: %v", q).Error(), "p")
		return
	}

	// A nil range represents to return all time-series value as response.
	if q.Range == nil {
		base = response{
			TimeInvolved: time.Since(stamp),
			Range: queryRange{
				Start: int64(math.MaxInt64),
				End:   int64(math.MinInt64),
			},
			Value: blockStream,
		}
		q.respond(encode(base))
		return
	}

	var (
		startPos, endPos = timeFirstBlock, timeLastBlock
		startPosFound    bool
	)
	for i := lengthBlockStream - 1; i >= 0; i-- {
		block := blockStream[i]
		if block.GetNormalizedTime() < q.Range.Start && !startPosFound {
			startPos = int64(i)
			startPosFound = true
			continue
		}
		if block.GetNormalizedTime() < q.Range.End && startPosFound {
			endPos = int64(i)
			break
		}
	}
	// we need to skip the first block in the range in order to fall inside
	// the range of [Start, End] query range. If not, this will lead to
	// Start - 1, End range of blocks.
	resultingBlockSlice := blockStream[endPos+1 : startPos+1]
	fmt.Println("resulting blocks are")
	fmt.Println(resultingBlockSlice)

	base = response{
		TimeInvolved: time.Since(stamp),
		Range:        *q.Range,
		Value:        resultingBlockSlice,
	}
	js = encode(base)
	q.respond(js)
}

// validate performs pre-validations on the query in order to avoid faults
// while traversing through the time-series values.
func (q *Query) validate(timeFirstBlock, timeLastBlock int64) bool {
	if q.Range.Start < timeLastBlock || q.Range.End > timeFirstBlock {
		q.ReturnNILResponse()
		return false
	}
	if q.Range.Start < timeLastBlock || q.Range.End > timeFirstBlock {
		q.ReturnNILResponse()
		return false
	}
	return true
}

// ReturnNILResponse is used only in cases when the response is known to be
// null value.
func (q *Query) ReturnNILResponse() {
	q.queryResponse.TimeInvolved = time.Since(q.stamp)
	q.respond(encode(q.queryResponse))
}

// ReturnMessageResponse is used only for sending simple messages. This should
// not be used to send JSON based results.
func (q *Query) ReturnMessageResponse(message string) {
	q.queryResponse = response{
		TimeInvolved: time.Since(q.stamp),
		Value:        message,
	}
	q.respond(encode(q.queryResponse))
}

func (q *Query) respond(js []byte) {
	if _, e := q.ws.Write(js); e != nil {
		logger.Terminal(fmt.Errorf("unable to respond: %s", e.Error()).Error(), "p")
	}
}
