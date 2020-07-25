package jitter

import (
	"strconv"
	"sync"
	"time"

	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"

	parser "github.com/zairza-cetb/bench-routes/src/lib/config"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
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
		logger.Terminal("invalid signal", "f")
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
				logger.Terminal("Not able to connect to external network please check you internet connection", "p")
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
			logger.Terminal("terminating jitter goroutine", "p")
			return
		default:
			logger.Terminal("invalid service-state value of jitter", "f")
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
			logger.Terminal("invalid interval-type for jitter", "f")
			return
		}
	}
}

func (ps *Jitter) jitter(urlRaw, machineHash string, packets int, wg *sync.WaitGroup) {
	resp, err := utils.CLIPing(urlRaw, packets)
	if err != nil {
		msg := "unable to reach " + urlRaw
		logger.Terminal(msg, "p")
		wg.Done()
		return
	}
	result := scrap.CLIJitterScrap(resp)
	newBlock := *tsdb.GetNewBlock("jitter", fToS(result))
	(*ps.targets)[machineHash].Jitter = (*ps.targets)[machineHash].Jitter.Append(newBlock)
	wg.Done()
}

func fToS(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
