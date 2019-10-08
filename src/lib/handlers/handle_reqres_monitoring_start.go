package handlers

import (
	"log"
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/modules/response"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
)

// HandleReqResMonitoringStart handle the route "start"
func HandleReqResMonitoringStart(config utils.YAMLBenchRoutesType, reqResMonitoringServiceState string) {
	routes := config.Config.Routes
	monitoringInterval := getInterval(config.Config.Interval, "req-res-delay-and-monitoring")
	if monitoringInterval == (testInterval{}) {
		log.Fatalf("interval not found in configuration file for req-res monitoring")
		return
	}

	for {
		config = config.Refresh()
		switch config.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring {
		case "active":
			var wg sync.WaitGroup
			wg.Add(len(routes))
			for _, route := range routes {
				go response.HandleResponseDelayForRoute(route, utils.GetHash(route.URL), &wg)
			}
			wg.Wait()
		case "passive":
			// terminate the goroutine
			log.Printf("terminating req-res monitoring goroutine\n")
			return
		default:
			log.Fatalf("invalid service-state value of req-res monitoring\n")
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
			log.Fatalf("invalid interval-type for req-res monitoring\n")
			return
		}
	}
}
