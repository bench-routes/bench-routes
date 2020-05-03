// contains the required struct types for the project in order to marshal and unmarshal
// the general tsdb block

package utils

import (
	"github.com/zairza-cetb/bench-routes/src/lib/parser"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

// Ping type for storing Ping values in TSDB
type Ping struct {
	Min  float64
	Mean float64
	Max  float64
	MDev float64
}

// PingResp type for storing Ping values in TSDB
type PingResp struct {
	Min            float64 `json:"Min"`
	Mean           float64 `json:"Mean"`
	Max            float64 `json:"Max"`
	MDev           float64 `json:"MDev"`
	NormalizedTime int64   `json:"NormalizedTime"`
	Timestamp      string  `json:"Timestamp"`
	Relative       int     `json:"relative"`
}

// JitterResp type for storing Ping values in TSDB
type JitterResp struct {
	Datapoint      float64 `json:"datapoint"`
	NormalizedTime int64   `json:"NormalizedTime"`
	Timestamp      string  `json:"Timestamp"`
	Relative       int     `json:"relative"`
}

// FloodPing type for storing Ping values in TSDB
type FloodPing struct {
	Min        float64
	Mean       float64
	Max        float64
	MDev       float64
	PacketLoss float64
}

// FloodPingResp type for storing Ping values in TSDB
type FloodPingResp struct {
	Min            float64 `json:"Min"`
	Mean           float64 `json:"Mean"`
	Max            float64 `json:"Max"`
	MDev           float64 `json:"MDev"`
	PacketLoss     float64 `json:"PacketLoss"`
	NormalizedTime int64   `json:"NormalizedTime"`
	Timestamp      string  `json:"Timestamp"`
	Relative       int     `json:"relative"`
}

// Response struct
// This is the object that we return from resp_delay module
// Contains delay in monitor and the monitor length
type Response struct {
	Delay         int
	ResLength     int
	ResStatusCode int
}

// ResponseResp for responding the querier.
type ResponseResp struct {
	Delay          int    `json:"delay"`
	ResLength      int    `json:"resLength"`
	ResStatusCode  int    `json:"resStatusCode"`
	NormalizedTime int64  `json:"NormalizedTime"`
	Timestamp      string `json:"Timestamp"`
	Relative       int    `json:"relative"`
}

// BRMatrix type for storing multi-dimensional information related to a route.
type BRMatrix struct {
	URL                                              string
	Domain                                           string
	Method                                           string
	Route                                            string
	Headers                                          []parser.Headers
	Params                                           []parser.Params
	PingChain, JitterChain, FPingChain, MonitorChain *tsdb.Chain
}

// BRmap forms a map of matrix that can be accessed with least possible time.
type BRmap map[string]*BRMatrix
