// contains the required struct types for the project in order to marshal and unmarshal
// the general tsdb block

package utils

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

// BlockPing block for ping case
type BlockPing struct {
	PrevBlock      *BlockPing
	NextBlock      *BlockPing
	Datapoint      Ping
	NormalizedTime int64
}

// BlockFloodPing block for ping case
type BlockFloodPing struct {
	PrevBlock      *BlockFloodPing
	NextBlock      *BlockFloodPing
	Datapoint      FloodPing
	NormalizedTime int64
}

// Response struct
// This is the object that we return from resp_delay module
// Contains delay in response and the response length
type Response struct {
	Delay         int
	ResLength     int64
	ResStatusCode int
}
