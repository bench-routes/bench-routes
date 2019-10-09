package controllers

import (
	core "github.com/zairza-cetb/bench-routes/src/lib"
)

// ReqResMonitoringController controlls requests for req-res
// monitoring and sends it to respective handler
func ReqResMonitoringController(signal string) bool {
	return core.HandleReqResGeneral(signal)
}
