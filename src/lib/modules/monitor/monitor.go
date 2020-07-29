package monitor

import (
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/common/log"
	parser "github.com/zairza-cetb/bench-routes/src/lib/config"
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/request"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/src/lib/utils/prom"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

// Monitor is the structure that implements the Monitoring service.
type Monitor struct {
	localConfig    *parser.Config
	scrapeInterval TestInterval
	targets        *map[string]*utils.BRMatrix
	test           bool
}

// TestInterval stores the value of the duration and the type of test
type TestInterval struct {
	OfType   string
	Duration int64
}

// New returns a Monitor type.
func New(configuration *parser.Config, scrapeInterval TestInterval, targets *map[string]*utils.BRMatrix) *Monitor {
	return &Monitor{
		localConfig:    configuration,
		scrapeInterval: scrapeInterval,
		targets:        targets,
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
		log.Errorln("invalid signal")
	}
	return false
}

// IsActive returns the current state of the service.
func (ps *Monitor) IsActive() bool {
	return ps.localConfig.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring == "active"
}

// perform carries out monitoring activities.
func (ps *Monitor) perform() {
	for {
		reqResMonitoringServiceState := ps.localConfig.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring
		monitoringInterval := getInterval(ps.localConfig.Config.Interval, "monitoring")
		if monitoringInterval == (TestInterval{}) {
			log.Errorln("interval not found in configuration file for req-res monitoring")
			return
		}
		switch reqResMonitoringServiceState {
		case "active":
			var wg sync.WaitGroup
			wg.Add(len(*ps.targets))
			for matrixHash, matrix := range *ps.targets {
				go ps.responseDelay(&wg, matrixHash, matrix.Route)
			}
			wg.Wait()
		case "passive":
			log.Infoln("terminating req-res monitoring goroutine")
			return
		default:
			log.Infoln("invalid service-state value of req-res monitoring")
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
			log.Infoln("invalid interval-type for req-res monitoring")
			return
		}
	}
}

func (ps *Monitor) responseDelay(wg *sync.WaitGroup, matrixHash string, route parser.Route) {
	response := make(chan request.ResponseWrapper)
	req := request.New(route.URL, request.ToMap(route.Header), request.ToMap(route.Params), request.ToMap(route.Body), route.Labels)
	stamp := time.Now()
	go req.Send(request.MethodUintPresentation(route.Method), response)
	resp := <-response
	var (
		delay      = time.Since(stamp).Milliseconds()
		resLength  = len(resp.ReponseStringified)
		statusCode = resp.Status
	)
	(*ps.targets)[matrixHash].Metrics.ResponseDelay.With(map[string]string{
		prom.LabelMethod: route.Method,
		prom.LabelDomain: filters.HTTPPingFilterValue(route.URL),
		prom.LabelURL:    route.URL,
	}).Set(float64(delay))
	(*ps.targets)[matrixHash].Metrics.ResponseLength.With(map[string]string{
		prom.LabelMethod: route.Method,
		prom.LabelDomain: filters.HTTPPingFilterValue(route.URL),
		prom.LabelURL:    route.URL,
	}).Set(float64(resLength))
	(*ps.targets)[matrixHash].Metrics.StatusCode.With(map[string]string{
		prom.LabelMethod: route.Method,
		prom.LabelDomain: filters.HTTPPingFilterValue(route.URL),
		prom.LabelURL:    route.URL,
	}).Set(float64(statusCode))
	(*ps.targets)[matrixHash].Metrics.MonitorCount.With(map[string]string{
		prom.LabelMethod: route.Method,
		prom.LabelDomain: filters.HTTPPingFilterValue(route.URL),
		prom.LabelURL:    route.URL,
	}).Inc()
	g := getNormalizedBlockString(utils.Response{
		Delay:         delay,
		ResLength:     len(resp.ReponseStringified),
		ResStatusCode: resp.Status,
	})
	block := *tsdb.GetNewBlock("req-res", g)
	(*ps.targets)[matrixHash].MonitorChain = (*ps.targets)[matrixHash].MonitorChain.Append(block)
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
	log.Infoln("interval not found in the configuration file")
	return TestInterval{}
}

// returns the stringified form of the combined data
func getNormalizedBlockString(b utils.Response) string {
	return strconv.Itoa(int(b.Delay)) + tsdb.BlockDataSeparator + strconv.Itoa(b.ResLength) + tsdb.BlockDataSeparator +
		strconv.Itoa(b.ResStatusCode)
}
