package utils

// TypePingScrap as datatype for ping outputs
type TypePingScrap struct {
	Min, Avg, Max, Mdev float64
}

// TypeFloodPingScrap as datatype for flood ping outputs
type TypeFloodPingScrap struct {
	Min, Avg, Max, Mdev, PacketLoss float64
}
