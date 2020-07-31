package ping

import (
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/common/log"
	parser "github.com/zairza-cetb/bench-routes/src/lib/config"
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/src/lib/utils/prom"
	"github.com/zairza-cetb/bench-routes/tsdb"
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

func (ps *Ping) perform(pingInterval TestInterval) {
	for {
		switch ps.localConfig.Config.UtilsConf.ServicesSignal.Ping {
		case "active":
			err, _ := utils.VerifyConnection()
			if !err {
				log.Warnln("unable to connect external network. please check you internet connection", "p")
			} else {
				var wg sync.WaitGroup
				wg.Add(len(*ps.targets))
				for machineHash, machineIP := range *ps.targets {
					go ps.ping(machineIP.IPDomain, machineHash, 3, &wg)
				}
				wg.Wait()
			}
		case "passive":
			// terminate the goroutine
			log.Infoln("terminating ping goroutine")
			return
		default:
			log.Warnln("invalid service-state value of ping")
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
			log.Infoln("invalid interval-type for ping")
			return
		}
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
