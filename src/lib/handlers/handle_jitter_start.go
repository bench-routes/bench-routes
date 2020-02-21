package handlers

import (
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/jitter"
	"github.com/zairza-cetb/bench-routes/src/lib/parser"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
)

// HandleJitterStart handle the route "start"
func HandleJitterStart(config *parser.YAMLBenchRoutesType, jitterServiceState string) {
	jitterConfig := config.Config.Routes
	jitterInterval := GetInterval(config.Config.Interval, "jitter")
	if jitterInterval == (TestInterval{}) {
		logger.Terminal("interval not found in configuration file for jitter", "f")
		return
	}
	urlStack := make(map[string]string)
	for _, route := range jitterConfig {
		url := route.URL
		urlHash := utils.GetHash(url)
		// maintain urls uniquely
		_, ok := urlStack[urlHash]
		if !ok {
			urlStack[urlHash] = *filters.HTTPPingFilter(&url)
		}
	}
	doJitter(config, urlStack, jitterInterval)
}

func doJitter(config *parser.YAMLBenchRoutesType, urlStack map[string]string, jitterInterval TestInterval) {
	i := 0
	for {
		i++
		config.Refresh()

		switch config.Config.UtilsConf.ServicesSignal.Jitter {
		case "active":
			var wg sync.WaitGroup
			wg.Add(len(urlStack))
			for _, u := range urlStack {
				go jitter.HandleJitter(utils.Jitterc, u, 10, u, &wg, false)
			}
			wg.Wait()
		case "passive":
			// terminate the goroutine
			logger.Terminal("terminating jitter goroutine", "p")
			return
		default:
			logger.Terminal("invalid service-state value of jitter", "p")
			return
		}

		intrv := time.Duration(jitterInterval.duration)
		switch jitterInterval.ofType {
		case "hr":
			time.Sleep(intrv * time.Hour)
		case "min":
			time.Sleep(intrv * time.Minute)
		case "sec":
			time.Sleep(intrv * time.Second)
		default:
			logger.Terminal("invalid interval-type for jitter", "f")
			return
		}
	}
}
