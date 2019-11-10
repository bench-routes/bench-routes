package main

import (
	"log"

	"github.com/zairza-cetb/bench-routes/src/lib/handlers"
)

type qPingRoute struct {
	URL string `json:"url"`
}

type qJitterRoute struct {
	URL string `json:"url"`
}

// HandlerPingGeneral handles the ping route
func HandlerPingGeneral(signal string) bool {

	// Get latest service state settings
	configuration = configuration.Refresh()
	pingServiceState := configuration.Config.UtilsConf.ServicesSignal.Ping

	switch signal {
	case "start":
		if pingServiceState == "passive" {

			configuration.Config.UtilsConf.ServicesSignal.Ping = "active"
			_, e := configuration.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				handlers.HandlePingStart(configuration, pingServiceState)
			}()
			return true
		}
		// return handlePingStart(configuration, pingServiceState)
	case "stop":
		configuration.Config.UtilsConf.ServicesSignal.Ping = "passive"
		_, e := configuration.Write()
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
	configuration = configuration.Refresh()
	jitterServiceState := configuration.Config.UtilsConf.ServicesSignal.Jitter

	switch signal {
	case "start":
		if jitterServiceState == "passive" {
			configuration.Config.UtilsConf.ServicesSignal.Jitter = "active"
			_, e := configuration.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				handlers.HandleJitterStart(configuration, jitterServiceState)
			}()
			return true
		}
	case "stop":
		configuration.Config.UtilsConf.ServicesSignal.Jitter = "passive"
		_, e := configuration.Write()
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
	configuration = configuration.Refresh()
	reqResMonitoringServiceState := configuration.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring

	switch signal {
	case "start":
		if reqResMonitoringServiceState == "passive" {

			configuration.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring = "active"
			_, e := configuration.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				handlers.HandleReqResMonitoringStart(configuration, reqResMonitoringServiceState)
			}()
			return true
		}
	case "stop":
		configuration.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring = "passive"
		_, e := configuration.Write()
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
	configuration = configuration.Refresh()
	serviceState := configuration.Config.UtilsConf.ServicesSignal.FloodPing

	switch signal {
	case "start":
		if serviceState == "passive" {
			configuration.Config.UtilsConf.ServicesSignal.FloodPing = "active"
			_, e := configuration.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				handlers.HandleFloodPingStart(configuration, serviceState)
			}()
			return true
		}
	case "stop":
		configuration.Config.UtilsConf.ServicesSignal.FloodPing = "passive"
		_, e := configuration.Write()
		if e != nil {
			panic(e)
		}
		return true
	default:
		log.Fatalf("invalid signal")
	}
	return false
}
