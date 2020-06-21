package monitor

import (
	"io/ioutil"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/logger"

	parser "github.com/zairza-cetb/bench-routes/src/lib/config"
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
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
		// if ps.isRunning {
		// 	ps.signalStop <- struct{}{}
		// }
		ps.localConfig.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring = "active"
		// ps.isRunning = true
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
		// select {
		// case <-ps.signalStop:
		// 	ps.isRunning = false
		// 	break
		// default:
		// 	ps.isRunning = true
		// }

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

	c := make(chan utils.Response)
	responseObject := routeDispatcher(route, c)

	g := getNormalizedBlockString(responseObject)
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

// routeDispatcher dispatches a route to respective handlers based on it's request type
func routeDispatcher(route parser.Route, c chan utils.Response) utils.Response {
	if route.Method == "GET" {
		return handleGetRequest(route.URL)
	}
	// If fail, then
	// send a very large integer to automatically rule out as it
	// is much much larger than the threshold
	return utils.Response{Delay: math.MaxInt32, ResLength: 0, ResStatusCode: 100}
}

// handleGetRequest specifically handles routes with GET Requests. Calculates timestamp before
// and after processing of each request and returns the difference
func handleGetRequest(url string) utils.Response {
	start := time.Now().UnixNano()
	resp := *utils.SendGETRequest(url)
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	resLength := len(content)
	respStatusCode := resp.StatusCode

	end := time.Now().UnixNano()
	diff := int((end - start) / int64(time.Millisecond))
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}

	return utils.Response{Delay: diff, ResLength: resLength, ResStatusCode: respStatusCode}
}

// returns the stringified form of the combined data
func getNormalizedBlockString(b utils.Response) string {
	return strconv.Itoa(b.Delay) + tsdb.BlockDataSeparator + strconv.Itoa(b.ResLength) + tsdb.BlockDataSeparator +
		strconv.Itoa(b.ResStatusCode)
}
