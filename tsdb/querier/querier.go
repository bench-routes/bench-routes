package querier

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

// Querier is a querying type.
type Querier struct {
	SocketInstance *websocket.Conn
	DBPath         string
	CollectorPath  string
}

// Query is a complex of db path, type and range of time-series.
type Query struct {
	ws            *websocket.Conn
	Path          string
	Type          string
	Range         *queryRange
	stamp         time.Time
	queryResponse response
}

// queryRange is the range of time within which the set of blocks
// are to be looked for.
type queryRange struct {
	start int64 `json:"start"`
	end   int64 `json:"end"`
}

// New returns a new Querier that implements querying.
func New(socketInstance *websocket.Conn, TSDBPath, collectorPath string) *Querier {
	return &Querier{
		SocketInstance: socketInstance,
		DBPath:         TSDBPath,
		CollectorPath:  collectorPath,
	}
}

// QueryBuilder builds a query that can be executed over a time-series
// range.
func (qs *Querier) QueryBuilder(ofType, middle, url string) *Query {
	path := fmt.Sprintf("%s/%s/chunk_%s/%s.json", qs.DBPath, ofType, middle, url)

	if ofType == "system-metrics" {
		path = fmt.Sprint("%s/system.json", ofType)
	}

	return &Query{
		ws:   qs.SocketInstance,
		Path: path,
		Type: ofType,
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
		start: fromTimestamp,
		end:   tillTimestamp,
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
		fmt.Errorf("invalid query received: ", q)
	}

	q.respond(js)
}

func (q *Query) validate(timeFirstBlock, timeLastBlock int64) bool {
	if q.Range.start < timeLastBlock || q.Range.end > timeFirstBlock {
		q.returnNILResponse()
		return false
	}
	if q.Range.start < timeLastBlock || q.Range.end > timeFirstBlock {
		q.returnNILResponse()
		return false
	}
	return true
}

// returnNILResponse is used only in cases when the response is known to be
// null value.
func (q *Query) returnNILResponse() {
	q.queryResponse.timeInvolved = time.Since(q.stamp)
	q.respond(encode(q.queryResponse))
}

// returnMessageResponse is used only for sending simple messages. This should
// not be used to send JSON based results.
func (q *Query) returnMessageResponse(message string) {
	q.queryResponse = response{
		timeInvolved: time.Since(q.stamp),
		value:        message,
	}
	q.respond(encode(q.queryResponse))
}

func (q *Query) respond(js []byte) {
	if e := q.ws.WriteMessage(1, js); e != nil {
		fmt.Errorf("unable to respond: %s", e.Error())
	}
}
