package handlers

import (
	"log"
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/ping"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

// HandlePingStart handle the route "start"
func HandlePingStart(config utils.YAMLBenchRoutesType, pingServiceState string) {
	pingConfig := config.Config.Routes
	pingInterval := getInterval(config.Config.Interval, "ping")
	if pingInterval == (testInterval{}) {
		log.Fatalf("interval not found in configuration file for ping")
		return
	}

	urlStack := make(map[string]string)
	for _, route := range pingConfig {
		url := route.URL
		urlHash := utils.GetHash(url)
		// maintain urls uniquely
		_, ok := urlStack[urlHash]
		if !ok {
			urlStack[urlHash] = *filters.HTTPPingFilter(&url)
		}
	}
	doPing(config, urlStack, pingInterval)
}

func doPing(config utils.YAMLBenchRoutesType, urlStack map[string]string, pingInterval testInterval) {
	i := 0
	for {
		i++
		config = config.Refresh()

		switch config.Config.UtilsConf.ServicesSignal.Ping {
		case "active":
			var wg sync.WaitGroup
			wg.Add(len(urlStack))
			for _, u := range urlStack {
				go ping.HandlePing(tsdb.GlobalPingChain, &u, 10, u, &wg, false)
			}

			wg.Wait()
		case "passive":
			// terminate the goroutine
			log.Printf("terminating ping goroutine\n")
			return
		default:
			log.Fatalf("invalid service-state value of ping\n")
			return
		}

		intrv := time.Duration(pingInterval.duration)
		switch pingInterval.ofType {
		case "hr":
			time.Sleep(intrv * time.Hour)
		case "min":
			time.Sleep(intrv * time.Minute)
		case "sec":
			time.Sleep(intrv * time.Second)
		default:
			log.Fatalf("invalid interval-type for ping\n")
			return
		}
	}
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
