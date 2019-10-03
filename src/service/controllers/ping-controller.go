package controllers

import (
	core "github.com/zairza-cetb/bench-routes/src/lib"
)

// PingController controllers the ping requests and transfers to the respective handler
func PingController(sig string) bool {
	return core.HandlerPingGeneral(sig)
}
