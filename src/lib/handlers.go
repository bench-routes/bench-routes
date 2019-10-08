package lib

import (
	"log"

	"github.com/zairza-cetb/bench-routes/src/lib/handlers"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
)

const (
	// ConfigurationFilePath is the constant path to the configuration file needed to start the application
	// written from root file since the application starts from `make run`
	ConfigurationFilePath = "storage/local-config.yml"
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
	// Configuration corresponds to the global use of configuration settings
	Configuration utils.YAMLBenchRoutesType
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

// HandlerFloodPingGeneral handles the flood ping module
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
