package tsdb

import (
	"github.com/gorilla/websocket"
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"sync"
)

const (
	directory = "storage/"
	prefix    = "chunk_"
	format    = ".json"
)

// BRQuerier is a querier for querying over the time-series stored in the tsdb
type BRQuerier struct {
	// As of now, ServiceName can be one of these: ping, food-ping, jitter, req-res-delay
	ServiceName string

	// DomainIP contains a unique identifier for uniquely identifying the particular time-series
	DomainIP string

	// Suffix is used for cases involving additional information for identifying the
	// db. Example: req-res-delay-monitoring contains `_delay`, `_length`, `_status`. Hence,
	// the suffix is the respective word excluding `_`
	Suffix string
	mux    sync.RWMutex

	// Reader (aka reader) reads the queries from the time-series samples and allows various
	// operations on it.
	reader Reader

	// Connection helps the querier to respond to the requests directly without
	// worrying about the goroutines.
	Connection *websocket.Conn
}

// FetchAllSeries returns all the series from the particular chunk. It is go-routine safe and
// respond to multiple same series requests.
func (bq *BRQuerier) FetchAllSeries() {
	bq.mux.RLock()
	defer bq.mux.RUnlock()

	switch bq.ServiceName {
	case "ping":
		series, err := bq.reader.open(bq.ServiceName, bq.fetchTSStorageAddress())
		if err != nil {
			panic(err)
		}
		if err := bq.Connection.WriteMessage(1, []byte(series)); err != nil {
			panic(err)
		}
	}
}

// Fetch time-series storage address
func (bq *BRQuerier) fetchTSStorageAddress() (address string) {
	switch bq.ServiceName {
	case "ping":
		address = directory + "ping/" + prefix + bq.ServiceName + "_" + *filters.HTTPPingFilter(&bq.DomainIP) + format
	case "jitter":
		address = directory + "jitter/" + prefix + bq.ServiceName + "_" + *filters.HTTPPingFilter(&bq.DomainIP) + format
	case "flood-ping":
		address = directory + "flood-ping/" + prefix + "flood_ping_" + *filters.HTTPPingFilter(&bq.DomainIP) + format
	case "req-res-delay":
		address = directory + "req-res-delay-monitoring/" + prefix + "req_res_" + filters.RouteDestroyer(bq.DomainIP) + format
	}
	return
}
