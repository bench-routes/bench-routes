package ping

import (
	"sync"
	"time"

	parser "github.com/zairza-cetb/bench-routes/src/lib/config"
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

// FloodPing is the structure that implements the Ping service.
type FloodPing struct {
	localConfig    *parser.Config
	scrapeInterval TestInterval
	targets        *map[string]*utils.MachineType
	password       string
	test           bool
}

// Newf returns a Flood Ping type.
func Newf(configuration *parser.Config, scrapeInterval TestInterval, password string, targets *map[string]*utils.MachineType) *FloodPing {
	return &FloodPing{
		localConfig:    configuration,
		scrapeInterval: scrapeInterval,
		password:       password,
		targets:        targets,
		test:           false,
	}
}

// FloodPingResponse is used to decode the tsdb blocks to data points that supports JSON encoding.
type FloodPingResponse struct {
	Min   string `json:"minValue"`
	Avg   string `json:"avgValue"`
	Max   string `json:"maxValue"`
	Mdev  string `json:"mdevValue"`
	Ploss string `json:"packetLoss"`
}

// Iteratef iterates over the local-configuration file to keep state
// of the ping service in sync with the local configuration.
// It is responsible for stopping the service without damaging the currently
// calculated samples.
func (ps *FloodPing) Iteratef(signal string, isTest bool) bool {
	if isTest {
		ps.test = true
	}
	switch signal {
	case "start":
		ps.localConfig.Config.UtilsConf.ServicesSignal.FloodPing = "active"
		go ps.setConfigurations()
		return true
	case "stop":
		ps.localConfig.Config.UtilsConf.ServicesSignal.FloodPing = "passive"
		return true
	default:
		logger.Terminal("invalid signal", "f")
	}
	return false
}

// GetServiceState returns the current state of the service.
func (ps *FloodPing) GetServiceState() bool {
	return ps.localConfig.Config.UtilsConf.ServicesSignal.FloodPing == "active"
}

func (ps *FloodPing) setConfigurations() {
	pingInterval := ps.scrapeInterval
	ps.perform(pingInterval)
}

func (ps *FloodPing) perform(pingInterval TestInterval) {
	for {
		switch ps.localConfig.Config.UtilsConf.ServicesSignal.FloodPing {
		case "active":
			err, _ := utils.VerifyConnection()
			if !err {
				logger.Terminal("Not able to connect to external network please check you internet connection", "p")
			} else {
				var wg sync.WaitGroup
				wg.Add(len(*ps.targets))
				for machineHash, machineIP := range *ps.targets {
					go ps.ping(machineIP.IPDomain, machineHash, 10, &wg)
				}
				wg.Wait()
			}
		case "passive":
			// terminate the goroutine
			logger.Terminal("terminating ping goroutine", "p")
			return
		default:
			logger.Terminal("invalid service-state value of ping", "f")
			return
		}

		intrv := time.Duration(*ps.localConfig.Config.Interval[0].Duration)
		switch pingInterval.OfType {
		case "hr":
			time.Sleep(intrv * time.Hour)
		case "min":
			time.Sleep(intrv * time.Minute)
		case "sec":
			time.Sleep(intrv * time.Second)
		default:
			logger.Terminal("invalid interval-type for ping", "f")
			return
		}
	}
}

func (ps *FloodPing) ping(urlRaw, machineHash string, packets int, wg *sync.WaitGroup) {
	password := ps.password
	resp, err := utils.CLIFloodPing(urlRaw, packets, password)
	if err != nil {
		logger.File(*resp, "p")
		wg.Done()
		return
	}
	result := scrap.CLIFLoodPingScrap(resp)
	newBlock := *tsdb.GetNewBlock("flood-ping", getNormalizedBlockStringFlood(*result))
	(*ps.targets)[machineHash].FPing = (*ps.targets)[machineHash].FPing.Append(newBlock)
	wg.Done()
}

func getNormalizedBlockStringFlood(v utils.TypeFloodPingScrap) string {
	return fToS(v.Min) + tsdb.BlockDataSeparator + fToS(v.Avg) + tsdb.BlockDataSeparator +
		fToS(v.Max) + tsdb.BlockDataSeparator + fToS(v.Mdev) + tsdb.BlockDataSeparator + fToS(v.PacketLoss)
}
