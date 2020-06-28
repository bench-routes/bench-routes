package monitor

import (
	"strconv"
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/logger"

	parser "github.com/zairza-cetb/bench-routes/src/lib/config"
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/request"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

const (
	// PathReqResDelay stores the default address of storage directory of ping data
	PathReqResDelay = "storage/req-res-delay-monitoring"
)

// Monitor is the structure that implements the Monitoring service.
type Monitor struct {
	localConfig    *parser.Config
	scrapeInterval TestInterval
	chain          *[]*tsdb.Chain
	test           bool
}

//TestInterval stores the value of the duration and the type of test
type TestInterval struct {
	OfType   string
	Duration int64
}

// New returns a Monitor type.
func New(configuration *parser.Config, scrapeInterval TestInterval, chain *[]*tsdb.Chain) *Monitor {
	return &Monitor{
		localConfig:    configuration,
		scrapeInterval: scrapeInterval,
		chain:          chain,
		test:           false,
	}
}

// Iterate iterates over the local-configuration file to keep state
// of the monitoring service in sync with the local configuration.
// It is responsible for stopping the service without damaging the currently
// calculated samples.
func (ps *Monitor) Iterate(signal string, isTest bool) bool {
	if isTest {
		ps.test = true
	}

	switch signal {
	case "start":
		ps.localConfig.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring = "active"
		go ps.perform()
		return true
	case "stop":
		ps.localConfig.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring = "passive"
		return true
	default:
		logger.Terminal("invalid signal", "f")
	}
	return false
}

// IsActive returns the current state of the service.
func (ps *Monitor) IsActive() bool {
	return ps.localConfig.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring == "active"
}

// perform carries out monitoring activities.
func (ps *Monitor) perform() {
	routes := ps.localConfig.Config.Routes

	for {
		reqResMonitoringServiceState := ps.localConfig.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring
		monitoringInterval := getInterval(ps.localConfig.Config.Interval, "req-res-delay-and-monitoring")
		if monitoringInterval == (TestInterval{}) {
			logger.Terminal("interval not found in configuration file for req-res monitoring", "f")
			return
		}
		switch reqResMonitoringServiceState {
		case "active":
			var wg sync.WaitGroup
			wg.Add(len(routes))
			// We send global chain arrays
			// of monitor delay, length and
			// statusCode in an array of type [][]*tsdb.Chain
			for _, route := range routes {
				go ps.responseDelay(&wg, route)
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

		intrv := time.Duration(*ps.localConfig.Config.Interval[2].Duration)
		switch monitoringInterval.OfType {
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

func (ps *Monitor) responseDelay(wg *sync.WaitGroup, route parser.Route) {
	responseChains := ps.chain
	routeSuffix := filters.RouteDestroyer(route.URL)
	// Init paths for request-monitor-monitoring
	path := PathReqResDelay + "/" + "chunk_req_res_" + routeSuffix + ".json"

	response := make(chan string)
	req := request.New(route.URL, request.ToMap(route.Header), request.ToMap(route.Params), request.ToMap(route.Body))
	// responseObject := routeDispatcher(route, c)
	stamp := time.Now()
	go req.Send(request.MethodUintPresentation(route.Method), response)
	resp := <-response

	g := getNormalizedBlockString(utils.Response{
		Delay:         time.Since(stamp).Seconds(),
		ResLength:     len(resp),
		ResStatusCode: 200,
	})
	block := *tsdb.GetNewBlock("req-res", g)

	for index := range *responseChains {
		if (*responseChains)[index].Path == path || ps.test {
			(*responseChains)[index] = (*responseChains)[index].Append(block)
			if ps.test {
				continue
			}
			break
		}
	}

	wg.Done()
}

// getInterval extracts the interval type and value from the config file for the particular test
func getInterval(intervals []parser.Interval, testName string) TestInterval {
	for _, intrv := range intervals {
		if testName == intrv.Test {
			return TestInterval{
				OfType:   intrv.Type,
				Duration: *intrv.Duration,
			}
		}
	}

	// if the execution reaches this then it is an error as the interval was not found in the configuration file
	logger.Terminal("interval not found in the configuration file", "pa")
	return TestInterval{}
}

// returns the stringified form of the combined data
func getNormalizedBlockString(b utils.Response) string {
	return strconv.FormatFloat(b.Delay, 'f', -1, 64) + tsdb.BlockDataSeparator + strconv.Itoa(b.ResLength) + tsdb.BlockDataSeparator +
		strconv.Itoa(b.ResStatusCode)
}
