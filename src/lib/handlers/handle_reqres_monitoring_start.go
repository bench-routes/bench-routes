package handlers

import (
	"github.com/zairza-cetb/bench-routes/src/lib/modules/monitor"
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/response"
	"github.com/zairza-cetb/bench-routes/src/lib/parser"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
)

// HandleReqResMonitoringStart handle the route "start"
func HandleReqResMonitoringStart(config *parser.YAMLBenchRoutesType, reqResMonitoringServiceState string) {
	routes := config.Config.Routes

	for {
		config.Refresh()
		monitoringInterval := GetInterval(config.Config.Interval, "req-res-delay-and-monitoring")
		if monitoringInterval == (TestInterval{}) {
			logger.Terminal("interval not found in configuration file for req-res monitoring", "f")
			return
		}
		switch config.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring {
		case "active":
			var wg sync.WaitGroup
			wg.Add(len(routes))
			// We send global chain arrays
			// of monitor delay, length and
			// statusCode in an array of type [][]*tsdb.Chain
			for _, route := range routes {
				go monitor.HandleResponseDelayForRoute(
					utils.RespMonitoringc,
					route,
					utils.GetHash(route.URL),
					&wg)
			}
			wg.Wait()
		case "passive":
			// terminate the goroutine
			logger.Terminal("terminating req-res monitoring goroutine", "p")
			return
		default:
			logger.Terminal("invalid service-state value of req-res monitoring", "f")
			return
		}

		intrv := time.Duration(monitoringInterval.duration)
		switch monitoringInterval.ofType {
		case "hr":
			time.Sleep(intrv * time.Hour)
		case "min":
			time.Sleep(intrv * time.Minute)
		case "sec":
			time.Sleep(intrv * time.Second)
		default:
			logger.Terminal("invalid interval-type for req-res monitoring", "f")
			return
		}
	}
}
