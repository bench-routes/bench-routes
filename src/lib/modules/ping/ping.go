package ping

import (
	"strconv"
	"sync"
	"time"

	parser "github.com/bench-routes/bench-routes/src/lib/config"
	scrap "github.com/bench-routes/bench-routes/src/lib/filters/scraps"
	"github.com/bench-routes/bench-routes/src/lib/utils"
	"github.com/bench-routes/bench-routes/src/lib/utils/prom"
	"github.com/bench-routes/bench-routes/tsdb"
	"github.com/prometheus/common/log"
)

// Ping is the structure that implements the Ping service.
type Ping struct {
	localConfig    *parser.Config
	targets        *map[string]*utils.MachineType
	scrapeInterval TestInterval
	test           bool
}

// TestInterval stores the value of the duration and the type of test
type TestInterval struct {
	OfType   string
	Duration int64
}

// Response is used to decode the tsdb blocks to data points that supports JSON encoding.
type Response struct {
	Min  string `json:"minValue"`
	Avg  string `json:"avgValue"`
	Max  string `json:"maxValue"`
	Mdev string `json:"mdevValue"`
}

// New returns a Ping type.
func New(configuration *parser.Config, scrapeInterval TestInterval, targets *map[string]*utils.MachineType) *Ping {
	return &Ping{
		localConfig:    configuration,
		scrapeInterval: scrapeInterval,
		test:           false,
		targets:        targets,
	}
}

// Iterate iterates over the local-configuration file to keep state
// of the ping service in sync with the local configuration.
// It is responsible for stopping the service without damaging the currently
// calculated samples.
func (ps *Ping) Iterate(signal string, isTest bool) bool {
	if isTest {
		ps.test = true
	}
	switch signal {
	case "start":
		ps.localConfig.Config.UtilsConf.ServicesSignal.Ping = "active"
		go ps.setConfigurations()
		return true
	case "stop":
		ps.localConfig.Config.UtilsConf.ServicesSignal.Ping = "passive"
		return true
	default:
		log.Errorln("invalid signal")
	}
	return false
}

// IsActive returns the current state of the service.
func (ps *Ping) IsActive() bool {
	return ps.localConfig.Config.UtilsConf.ServicesSignal.Ping == "active"
}

func (ps *Ping) setConfigurations() {
	pingInterval := ps.scrapeInterval
	ps.perform(pingInterval)
}

// pingerMap holds the closing channel for current set of pingable machines workers
type pingerMap map[string]chan struct{}

// pinger is the worker, which pings an endpoint on regular intervals until the channels in closed
func (ps *Ping) pinger(urlRaw, machineHash string, packets int, intrv time.Duration, closeCh <-chan struct{}) {
	var wg sync.WaitGroup
	for { // This will infinitely run the ping as long as the channeled is not closed
		select {
		case <-closeCh: // Close channel
			return
		default:
			wg.Add(1)
			ps.ping(urlRaw, machineHash, 3, &wg)

		}

		time.Sleep(intrv) // Taking a nap
	}

}

func (ps *Ping) perform(pingInterval TestInterval) {

	var pmap = make(pingerMap)
	for {
		switch ps.localConfig.Config.UtilsConf.ServicesSignal.Ping {
		case "active":
			err, _ := utils.VerifyConnection()
			if !err {
				log.Warnln("unable to connect external network. please check you internet connection", "p")
			} else {
				interval := getPingInterval(time.Duration(*ps.localConfig.Config.Interval[0].Duration),
					pingInterval.OfType)
				if interval == 0 {
					return // Invalid interval unit
				}
				for machineHash, machineIP := range *ps.targets {
					_, ok := pmap[machineHash]
					if !ok { // We will spin a new worker, not found one
						ch := make(chan struct{})
						pmap[machineHash] = ch
						go ps.pinger(machineIP.IPDomain, // Registering a new go-routine
							machineHash, 3,
							interval, ch)

					}
				}
				// Closing unecessary go-routines
				for k, v := range pmap {
					_, ok := (*ps.targets)[k]
					if !ok {
						close(v)
						delete(pmap, k)
					}
				}

			}
		case "passive":
			// terminate the goroutine
			log.Infoln("terminating ping goroutine")
			return
		default:
			log.Warnln("invalid service-state value of ping")
			return
		}

	}
}

func getPingInterval(intrv time.Duration, ofType string) time.Duration {
	switch ofType {
	case "hr":
		return intrv * time.Hour
	case "min":
		return intrv * time.Minute
	case "sec":
		return intrv * time.Second
	default:
		log.Infoln("invalid interval-type for ping")
		return 0
	}

}

func (ps *Ping) ping(urlRaw, machineHash string, packets int, wg *sync.WaitGroup) {
	resp, err := utils.CLIPing(urlRaw, packets)
	if err != nil {
		log.Warnln("unable to reach " + urlRaw)
		wg.Done()
		return
	}
	result := scrap.CLIPingScrap(resp)
	(*ps.targets)[machineHash].Metrics.Ping.With(map[string]string{
		prom.LabelDomain:    urlRaw,
		prom.LabelPingTypes: "min",
	}).Set(result.Min)
	(*ps.targets)[machineHash].Metrics.Ping.With(map[string]string{
		prom.LabelDomain:    urlRaw,
		prom.LabelPingTypes: "mean",
	}).Set(result.Avg)
	(*ps.targets)[machineHash].Metrics.Ping.With(map[string]string{
		prom.LabelDomain:    urlRaw,
		prom.LabelPingTypes: "max",
	}).Set(result.Max)
	(*ps.targets)[machineHash].Metrics.Ping.With(map[string]string{
		prom.LabelDomain:    urlRaw,
		prom.LabelPingTypes: "mdev",
	}).Set(result.Mdev)
	(*ps.targets)[machineHash].Metrics.PingCount.With(map[string]string{
		prom.LabelDomain: urlRaw,
	}).Inc()
	newBlock := *tsdb.GetNewBlock("ping", getNormalizedBlockString(*result))
	(*ps.targets)[machineHash].Ping = (*ps.targets)[machineHash].Ping.Append(newBlock)
	wg.Done()
}

func getNormalizedBlockString(v utils.TypePingScrap) string {
	return fToS(v.Min) + tsdb.BlockDataSeparator + fToS(v.Avg) +
		tsdb.BlockDataSeparator + fToS(v.Max) + tsdb.BlockDataSeparator + fToS(v.Mdev)
}

func fToS(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
