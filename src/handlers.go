package main

import (
	"github.com/zairza-cetb/bench-routes/src/lib/handlers"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
)

type qPingRoute struct {
	URL string `json:"url"`
}

type qFloodPingRoute struct {
	URL string `json:"url"`
}

type qJitterRoute struct {
	URL string `json:"url"`
}

type qReqResDelayRoute struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

// HandlerPingGeneral handles the ping route
func HandlerPingGeneral(signal string) bool {

	// Get latest service state settings
	conf.Refresh()
	pingServiceState := conf.Config.UtilsConf.ServicesSignal.Ping

	switch signal {
	case "start":
		if pingServiceState == "passive" {

			conf.Config.UtilsConf.ServicesSignal.Ping = "active"
			_, e := conf.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				handlers.HandlePingStart(conf, pingServiceState)
			}()
			return true
		}
		// return handlePingStart(conf, pingServiceState)
	case "stop":
		conf.Config.UtilsConf.ServicesSignal.Ping = "passive"
		_, e := conf.Write()
		if e != nil {
			panic(e)
		}
		return true
	default:
		logger.Terminal("invalid signal", "f")
	}
	return false
}

//HandlerJitterGeneral handles the request to the jitter module
func HandlerJitterGeneral(signal string) bool {

	// Get latest service state settings
	conf.Refresh()
	jitterServiceState := conf.Config.UtilsConf.ServicesSignal.Jitter

	switch signal {
	case "start":
		if jitterServiceState == "passive" {
			conf.Config.UtilsConf.ServicesSignal.Jitter = "active"
			_, e := conf.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				handlers.HandleJitterStart(conf, jitterServiceState)
			}()
			return true
		}
	case "stop":
		conf.Config.UtilsConf.ServicesSignal.Jitter = "passive"
		_, e := conf.Write()
		if e != nil {
			panic(e)
		}
		return true
	default:
		logger.Terminal("invalid signal", "f")
	}
	return false
}

// HandleReqResGeneral is the handler for requests regarding
// req-res delay and monitoring
func HandleReqResGeneral(signal string) bool {
	// Get latest service state settings
	conf.Refresh()
	reqResMonitoringServiceState := conf.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring

	switch signal {
	case "start":
		if reqResMonitoringServiceState == "passive" {

			conf.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring = "active"
			_, e := conf.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				handlers.HandleReqResMonitoringStart(conf, reqResMonitoringServiceState)
			}()
			return true
		}
	case "stop":
		conf.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring = "passive"
		_, e := conf.Write()
		if e != nil {
			panic(e)
		}
		return true
	default:
		logger.Terminal("invalid signal", "f")
	}
	return false
}

//HandlerFloodPingGeneral handles the flood-ping route
func HandlerFloodPingGeneral(signal string) bool {
	// Refresh conf with latest update
	conf.Refresh()
	serviceState := conf.Config.UtilsConf.ServicesSignal.FloodPing

	switch signal {
	case "start":
		if serviceState == "passive" {
			conf.Config.UtilsConf.ServicesSignal.FloodPing = "active"
			_, e := conf.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				handlers.HandleFloodPingStart(conf, serviceState)
			}()
			return true
		}
	case "stop":
		conf.Config.UtilsConf.ServicesSignal.FloodPing = "passive"
		_, e := conf.Write()
		if e != nil {
			panic(e)
		}
		return true
	default:
		logger.Terminal("invalid signal", "f")
	}
	return false
}
