package querier

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	sysMetrics "github.com/zairza-cetb/bench-routes/src/metrics/system"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

// Querier is a querying type.
type Querier struct {
	SocketInstance *websocket.Conn
	DBPath         string
	CollectorPath  string
}

type Query struct {
	ws    *websocket.Conn
	Path  string
	Type  string
	Range string
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
		// range "" represents all time-series as querying
		Range: "",
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
	switch q.Type {
	case "system-metrics":
		var responses []sysMetrics.Response
		for _, b := range blockStream {
			responses = append(responses, sysMetrics.Decode(b.Datapoint))
		}
		js, e := json.Marshal(responses)
		if e != nil {
			panic(e)
		}
	}

	if e := q.ws.WriteMessage(1, js); e != nil {
		panic(e)
	}
}
