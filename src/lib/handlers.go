package lib

import (
	"log"

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
		return handlePingStart(Configuration, pingServiceState)
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

type testInterval struct {
	ofType   string
	duration int64
}

func getInterval(intervals []utils.Interval, testName string) testInterval {
	for _, intrv := range intervals {
		if testName == intrv.Test {
			return testInterval{
				ofType:   intrv.Type,
				duration: intrv.Duration,
			}
		}
	}

	// if the execution reaches this then it is an error as the interval was not found in the configuration file
	log.Panicf("interval not found in the configuration file\n")
	return testInterval{}
}
