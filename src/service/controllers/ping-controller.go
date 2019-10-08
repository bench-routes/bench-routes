package controllers

import (
	core "github.com/zairza-cetb/bench-routes/src/lib"
)

// PingController controllers the ping requests and transfers to the respective handler
func PingController(sig string) bool {
	return core.HandlerPingGeneral(sig)
}

// FloodPingController controls the flood ping requests and transfers it to handlers
func FloodPingController(sig string) bool {
	return core.HandlerFloodPingGeneral(sig)
}
