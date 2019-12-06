package tsdb

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
)

const (
	directory = "storage/"
	prefix    = "chunk_"
	format    = ".json"
)

// BQRoute is a struct which contains a Domain IP and it's subsequent method
type BQRoute struct {
	DomainIP string
	Method   string
}

// BRQuerier is a querier for querying over the time-series stored in the tsdb
type BRQuerier struct {
	// As of now, ServiceName can be one of these: ping, food-ping, jitter, req-res-delay
	ServiceName string

	// Route contains a unique identifier and it's method for uniquely
	// identifying the particular time-series
	Route BQRoute

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

	// Complete address of the persistent blocks.
	Address string
}

// FetchAllSeries returns all the series from the particular chunk.
// It does not do any processing, rather just a plain simple fetch and return the fetched samples.
// It is go-routine safe and respond to multiple same series requests.
func (bq *BRQuerier) FetchAllSeries() {
	bq.mux.RLock()
	defer bq.mux.RUnlock()

	switch bq.ServiceName {
	case "ping":
		series, err := bq.reader.open(bq.fetchTSStorageAddress())
		if err != nil {
			panic(err)
		}
		if err := bq.Connection.WriteMessage(1, []byte(series)); err != nil {
			panic(err)
		}
	case "flood-ping":
		series, err := bq.reader.open(bq.fetchTSStorageAddress())
		if err != nil {
			panic(err)
		}
		if err := bq.Connection.WriteMessage(1, []byte(series)); err != nil {
			panic(err)
		}
	case "jitter":
		series, err := bq.reader.open(bq.fetchTSStorageAddress())
		if err != nil {
			panic(err)
		}
		if err := bq.Connection.WriteMessage(1, []byte(series)); err != nil {
			panic(err)
		}
	case "req-res-delay":
		series, err := bq.reader.open(bq.fetchTSStorageAddress())
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
		address = directory + "ping/" + prefix + bq.ServiceName + "_" + *filters.HTTPPingFilter(&bq.Route.DomainIP) + format
	case "jitter":
		address = directory + "jitter/" + prefix + bq.ServiceName + "_" + *filters.HTTPPingFilter(&bq.Route.DomainIP) + format
	case "flood-ping":
		address = directory + "flood-ping/" + prefix + "flood_ping_" + *filters.HTTPPingFilter(&bq.Route.DomainIP) + format
	case "req-res-delay":
		address = directory + "req-res-delay-monitoring/" + prefix + "req_res_" + filters.RouteDestroyer(bq.Route.DomainIP) + format
	}
	bq.Address = address
	return
}

// GetAddress returns the address of the persistent blocks in the db.
func (bq *BRQuerier) GetAddress() string {
	return bq.Address
}
