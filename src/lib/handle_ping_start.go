package lib

import (
	"log"
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/ping"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

func handlePingStart(config utils.YAMLBenchRoutesType, pingServiceState string) bool {
	if pingServiceState == "passive" {

		Configuration.Config.UtilsConf.ServicesSignal.Ping = "active"
		_, e := Configuration.Write()
		if e != nil {
			panic(e)
		}
		go func() {
			pingConfig := Configuration.Config.Routes
			pingInterval := getInterval(Configuration.Config.Interval, "ping")
			if pingInterval == (testInterval{}) {
				log.Fatalf("interval not found in configuration file for ping")
				return
			}

			urlStack := make(map[string]string)
			urlHashMap := make(map[string]string)
			for _, route := range pingConfig {
				url := route.URL
				urlHash := utils.GetHash(url)
				// maintain urls uniquely
				_, ok := urlStack[urlHash]
				if !ok {
					urlStack[urlHash] = *filters.HTTPPingFilter(&url)
					urlHashMap[url] = urlHash
				}
			}
			doPing(urlStack, urlHashMap, pingInterval)
		}()
		return true
	}
	return false
}

func doPing(urlStack, urlHashMap map[string]string, pingInterval testInterval) {
	i := 0
	for {
		i++
		Configuration = Configuration.Refresh()

		switch Configuration.Config.UtilsConf.ServicesSignal.Ping {
		case "active":
			var wg sync.WaitGroup
			wg.Add(len(urlStack))
			for _, u := range urlStack {
				go ping.HandlePing(tsdb.GlobalPingChain, &u, 10, urlHashMap[urlStack[u]], &wg, false)
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
