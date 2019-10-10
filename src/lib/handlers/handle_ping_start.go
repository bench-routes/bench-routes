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
	pingInterval := GetInterval(config.Config.Interval, "ping")
	if pingInterval == (TestInterval{}) {
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

func doPing(config utils.YAMLBenchRoutesType, urlStack map[string]string, pingInterval TestInterval) {
	i := 0
	for {
		i++
		config = config.Refresh()

		switch config.Config.UtilsConf.ServicesSignal.Ping {
		case "active":
			var wg sync.WaitGroup
			wg.Add(len(urlStack))
			for _, u := range urlStack {
				go ping.HandlePing(tsdb.GlobalPingChain, u, 10, u, &wg, false)
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

// HandleFloodPingStart starts the flood ping service
func HandleFloodPingStart(config utils.YAMLBenchRoutesType, floodPingServiceState string) {
	floodPingConfig := config.Config.Routes
	floodPingInterval := GetInterval(config.Config.Interval, "ping")
	if floodPingInterval == (TestInterval{}) {
		log.Fatalf("interval not found in configuration file for ping")
		return
	}

	urlStack := make(map[string]string)
	for _, route := range floodPingConfig {
		url := route.URL
		urlHash := utils.GetHash(url)
		// maintain urls uniquely
		_, ok := urlStack[urlHash]
		if !ok {
			urlStack[urlHash] = *filters.HTTPPingFilter(&url)
		}
	}

	doFloodPing(config, urlStack, floodPingInterval)
}

func doFloodPing(config utils.YAMLBenchRoutesType, urlStack map[string]string, interval TestInterval) {
	i := 0
	for {
		i++
		config = config.Refresh()

		switch config.Config.UtilsConf.ServicesSignal.FloodPing {
		case "active":
			var wg sync.WaitGroup
			wg.Add(len(urlStack))
			for _, u := range urlStack {
				go ping.HandleFloodPing(tsdb.GlobalFloodPingChain, u, 10000, u, &wg, false, config.Config.Password)
			}

			wg.Wait()
		case "passive":
			//terminate the goroutine
			log.Printf("terminating flood ping goroutine\n")
			return
		default:
			log.Fatalf("invalid service-state value for flood-ping\n")
			return
		}

		intrv := time.Duration(interval.duration)
		switch interval.ofType {
		case "hr":
			time.Sleep(intrv * time.Hour)
		case "min":
			time.Sleep(intrv * time.Minute)
		case "sec":
			time.Sleep(intrv * time.Second)
		default:
			log.Fatalf("invalid interval-type for flood-ping\n")
			return
		}
	}
}

//TestInterval stores the value of the duration and the type of test
type TestInterval struct {
	ofType   string
	duration int64
}

//GetInterval extracts the interval type and value from the config file for the particular test
func GetInterval(intervals []utils.Interval, testName string) TestInterval {
	for _, intrv := range intervals {
		if testName == intrv.Test {
			return TestInterval{
				ofType:   intrv.Type,
				duration: intrv.Duration,
			}
		}
	}

	// if the execution reaches this then it is an error as the interval was not found in the configuration file
	log.Panicf("interval not found in the configuration file\n")
	return TestInterval{}
}
