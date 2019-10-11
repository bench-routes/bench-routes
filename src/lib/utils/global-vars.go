package utils

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

// TypePingScrap as datatype for ping outputs
type TypePingScrap struct {
	Min, Avg, Max, Mdev float64
}

// TypeFloodPingScrap as datatype for flood ping outputs
type TypeFloodPingScrap struct {
	Min, Avg, Max, Mdev, PacketLoss float64
}

// Response struct
// This is the object that we return from resp_delay module
// Contains delay in response and the response length
type Response struct {
	Delay         int
	ResLength     int64
	ResStatusCode int
}
