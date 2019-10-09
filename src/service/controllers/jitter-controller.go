package controllers

import (
	core "github.com/zairza-cetb/bench-routes/src/lib"
)

// JitterController controllers the jitter requests and transfers to the respective handler
func JitterController(sig string) bool {
	return core.HandlerJitterGeneral(sig)
}
