package utils

import (
	"github.com/zairza-cetb/bench-routes/tsdb"
)

const (
	// ConfigurationFilePath is the constant path to the configuration file needed to start the application
	// written from root file since the application starts from `make run`
	ConfigurationFilePath = "local-config.yml"
	// PathPing stores the defualt address of storage directory of ping data
	PathPing = "storage/ping"
	// PathJitter stores the defualt address of storage directory of jitter data
	PathJitter = "storage/jitter"
	// PathFloodPing stores the defualt address of storage directory of flood ping data
	PathFloodPing = "storage/flood-ping"
	// PathReqResDelayMonitoring stores the defualt address of storage directory of req-res and monitoring data
	PathReqResDelayMonitoring = "storage/req-res-delay-monitoring"
)

var (
	// Pingc contains chains of all the pings operating in bench-routes which has to be globally accessed
	// This is necessary as it helps to retain the parent values which are required for concurreny
	Pingc []*tsdb.Chain

	// Jitterc is the global chain array which can be used to maintain a list of chains that represent
	// the time-series values
	Jitterc []*tsdb.Chain

	// FPingc contains chains of flood ping operations in bench-routes which has to be globally accessed
	// This is necessary as it helps to retain the parent values which are required for concurreny
	FPingc []*tsdb.Chain

	// RespMonitoringc stands for Request-Response-Delay
	RespMonitoringc []*tsdb.Chain
)

// TypePingScrap as datatype for ping outputs
type TypePingScrap struct {
	Min, Avg, Max, Mdev float64
}

// TypeFloodPingScrap as datatype for flood ping outputs
type TypeFloodPingScrap struct {
	Min, Avg, Max, Mdev, PacketLoss float64
}
