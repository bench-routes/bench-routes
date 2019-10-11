package lib

import (
	"log"

	"github.com/zairza-cetb/bench-routes/src/lib/handlers"
	"github.com/zairza-cetb/bench-routes/src/lib/utils/parser"
)

var (
	// Configuration corresponds to the global use of configuration settings
	Configuration parser.YAMLBenchRoutesType
)

// HandlerPingGeneral handles the ping route
func HandlerPingGeneral(signal string) bool {

	// Get latest service state settings
	Configuration = Configuration.Refresh()
	pingServiceState := Configuration.Config.UtilsConf.ServicesSignal.Ping

	switch signal {
	case "start":
		if pingServiceState == "passive" {

			Configuration.Config.UtilsConf.ServicesSignal.Ping = "active"
			_, e := Configuration.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				handlers.HandlePingStart(Configuration, pingServiceState)
			}()
			return true
		}
		// return handlePingStart(Configuration, pingServiceState)
	case "stop":
		Configuration.Config.UtilsConf.ServicesSignal.Ping = "passive"
		_, e := Configuration.Write()
		if e != nil {
			panic(e)
		}
		return true
	default:
		log.Fatalf("invalid signal")
	}
	return false
}

//HandlerJitterGeneral handles the request to the jitter module
func HandlerJitterGeneral(signal string) bool {

	// Get latest service state settings
	Configuration = Configuration.Refresh()
	jitterServiceState := Configuration.Config.UtilsConf.ServicesSignal.Jitter

	switch signal {
	case "start":
		if jitterServiceState == "passive" {
			Configuration.Config.UtilsConf.ServicesSignal.Jitter = "active"
			_, e := Configuration.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				handlers.HandleJitterStart(Configuration, jitterServiceState)
			}()
			return true
		}
	case "stop":
		Configuration.Config.UtilsConf.ServicesSignal.Jitter = "passive"
		_, e := Configuration.Write()
		if e != nil {
			panic(e)
		}
		return true
	default:
		log.Fatalf("invalid signal")
	}
	return false
}

// HandleReqResGeneral is the handler for requests regarding
// req-res delay and monitoring
func HandleReqResGeneral(signal string) bool {
	// Get latest service state settings
	Configuration = Configuration.Refresh()
	reqResMonitoringServiceState := Configuration.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring

	switch signal {
	case "start":
		if reqResMonitoringServiceState == "passive" {

			Configuration.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring = "active"
			_, e := Configuration.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				handlers.HandleReqResMonitoringStart(Configuration, reqResMonitoringServiceState)
			}()
			return true
		}
	case "stop":
		Configuration.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring = "passive"
		_, e := Configuration.Write()
		if e != nil {
			panic(e)
		}
		return true
	default:
		log.Fatalf("invalid signal")
	}
	return false
}

//HandlerFloodPingGeneral handles the flood-ping route
func HandlerFloodPingGeneral(signal string) bool {
	// Refresh configuration with latest update
	Configuration = Configuration.Refresh()
	serviceState := Configuration.Config.UtilsConf.ServicesSignal.FloodPing

	switch signal {
	case "start":
		if serviceState == "passive" {
			Configuration.Config.UtilsConf.ServicesSignal.FloodPing = "active"
			_, e := Configuration.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				handlers.HandleFloodPingStart(Configuration, serviceState)
			}()
			return true
		}
	case "stop":
		Configuration.Config.UtilsConf.ServicesSignal.FloodPing = "passive"
		_, e := Configuration.Write()
		if e != nil {
			panic(e)
		}
		return true
	default:
		log.Fatalf("invalid signal")
	}
	return false
}
