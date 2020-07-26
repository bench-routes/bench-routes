// contains the required struct types for the project in order to marshal and unmarshal
// the general tsdb block

package utils

import (
	"github.com/prometheus/client_golang/prometheus"
	// parser "github.com/zairza-cetb/bench-routes/src/lib/config"
	parser "github.com/zairza-cetb/bench-routes/src/lib/config"
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
	Delay         int64 `json:"delay"`
	ResLength     int   `json:"resLength"`
	ResStatusCode int   `json:"resStatusCode"`
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

// ResponseTSDBChains as response type for information on the tsdb chains.
type ResponseTSDBChains struct {
	Name string    `json:"name"`
	Path ChainPath `json:"path"`
}

// ChainPath keeps the path of all chains corresponding to a matrix value.
type ChainPath struct {
	InstanceKey string `json:"matrixName"`
	Ping        string `json:"ping"`
	Jitter      string `json:"jitter"`
	Fping       string `json:"fping"`
	Monitor     string `json:"monitor"`
}

// MatrixResponse wraps the block stream of the chains.
type MatrixResponse struct {
	PingBlocks    []byte `json:"ping"`
	JitterBlocks  []byte `json:"jitter"`
	FpingBlocks   []byte `json:"fping"`
	MonitorBlocks []byte `json:"monitor"`
}

// BRMatrix type for storing multi-dimensional information related to a route.
type BRMatrix struct {
	FullURL                  string
	Route                    parser.Route
	PingChain, JitterChain   *tsdb.Chain
	FPingChain, MonitorChain *tsdb.Chain
	Metrics                  *EndpointMetrics
}

// MachineType type for targets. These are servers/vm-instances/load-balancer
// that are accessible through pure IP/domain address.
type MachineType struct {
	IPDomain string
	Ping     *tsdb.Chain
	Jitter   *tsdb.Chain
	FPing    *tsdb.Chain
	Metrics  *MachineMetrics
}

// MachineMetrics represent the prometheus metrics for each
// target machine/IP/load-balancer.
type MachineMetrics struct {
	Ping        *prometheus.GaugeVec
	PingCount   *prometheus.CounterVec
	Jitter      *prometheus.GaugeVec
	JitterCount *prometheus.CounterVec
	FPing       *prometheus.GaugeVec
	FPingCount  *prometheus.CounterVec
}

// EndpointMetrics represent the metrics based on API endpoints.
type EndpointMetrics struct {
	ResponseLength *prometheus.GaugeVec
	ResponseDelay  *prometheus.GaugeVec
	StatusCode     *prometheus.GaugeVec
	MonitorCount   *prometheus.CounterVec
}
