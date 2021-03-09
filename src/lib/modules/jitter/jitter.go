package jitter

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

// Jitter is the structure that implements the Jitter service.
type Jitter struct {
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
	Value string `json:"value"`
}

// New returns a Jitter type.
func New(configuration *parser.Config, scrapeInterval TestInterval, targets *map[string]*utils.MachineType) *Jitter {
	return &Jitter{
		localConfig:    configuration,
		scrapeInterval: scrapeInterval,
		targets:        targets,
		test:           false,
	}
}

// Iterate iterates over the local-configuration file to keep state
// of the jitter service in sync with the local configuration.
// It is responsible for stopping the service without damaging the currently
// calculated samples.
func (ps *Jitter) Iterate(signal string, isTest bool) bool {
	if isTest {
		ps.test = true
	}
	switch signal {
	case "start":
		ps.localConfig.Config.UtilsConf.ServicesSignal.Jitter = "active"
		go ps.setConfigurations()
		return true
	case "stop":
		ps.localConfig.Config.UtilsConf.ServicesSignal.Jitter = "passive"
		return true
	default:
		log.Errorln("invalid signal")
	}
	return false
}

// IsActive returns the current state of the service.
func (ps *Jitter) IsActive() bool {
	return ps.localConfig.Config.UtilsConf.ServicesSignal.Jitter == "active"
}

func (ps *Jitter) setConfigurations() {
	interval := ps.scrapeInterval
	ps.perform(interval)
}

func (ps *Jitter) perform(pingInterval TestInterval) {
	for {
		switch ps.localConfig.Config.UtilsConf.ServicesSignal.Jitter {
		case "active":
			err, _ := utils.VerifyConnection()
			if !err {
				log.Warnln("Not able to connect to external network please check you internet connection")
			} else {
				var wg sync.WaitGroup
				wg.Add(len(*ps.targets))
				for machineHash, machineIP := range *ps.targets {
					go ps.jitter(machineIP.IPDomain, machineHash, 3, &wg)
				}
				wg.Wait()
			}
		case "passive":
			// terminate the goroutine
			log.Infoln("terminating jitter goroutine")
			return
		default:
			log.Errorln("invalid service-state value of jitter")
			return
		}

		intrv := time.Duration(*ps.localConfig.Config.Interval[1].Duration)
		switch pingInterval.OfType {
		case "hr":
			time.Sleep(intrv * time.Hour)
		case "min":
			time.Sleep(intrv * time.Minute)
		case "sec":
			time.Sleep(intrv * time.Second)
		default:
			log.Errorln("invalid interval-type for jitter")
			return
		}
	}
}

func (ps *Jitter) jitter(urlRaw, machineHash string, packets int, wg *sync.WaitGroup) {
	resp, err := utils.CLIPing(urlRaw, packets)
	if err != nil {
		log.Warnln("unable to reach " + urlRaw)
		wg.Done()
		return
	}
	result := scrap.CLIJitterScrap(resp)
	(*ps.targets)[machineHash].Metrics.Jitter.With(map[string]string{
		prom.LabelDomain: urlRaw,
	}).Set(result)
	(*ps.targets)[machineHash].Metrics.JitterCount.With(map[string]string{
		prom.LabelDomain: urlRaw,
	}).Inc()
	newBlock := *tsdb.GetNewBlock("jitter", fToS(result))
	(*ps.targets)[machineHash].Jitter = (*ps.targets)[machineHash].Jitter.Append(newBlock)
	wg.Done()
}

func fToS(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
