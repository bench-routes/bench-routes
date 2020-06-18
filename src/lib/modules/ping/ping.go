package ping

import (
	"strconv"
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/config"
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

// Ping is the structure that implements the Ping service.
type Ping struct {
	localConfig    *parser.Config
	scrapeInterval TestInterval
	chain          *[]*tsdb.Chain
	test           bool
	isRunning      bool
	// signalStop is mostly used to reload the configurations of
	// new urls.
	signalStop chan struct{}
}

//TestInterval stores the value of the duration and the type of test
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
func New(configuration *parser.Config, scrapeInterval TestInterval, chain *[]*tsdb.Chain) *Ping {
	return &Ping{
		localConfig:    configuration,
		scrapeInterval: scrapeInterval,
		chain:          chain,
		test:           false,
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

	conf := ps.localConfig
	conf.Refresh()
	pingServiceState := conf.Config.UtilsConf.ServicesSignal.Ping

	switch signal {
	case "start":
		if pingServiceState == "passive" {
			if ps.isRunning {
				ps.signalStop <- struct{}{}
			}
			conf.Config.UtilsConf.ServicesSignal.Ping = "active"
			_, e := conf.Write()
			if e != nil {
				panic(e)
			}
			ps.isRunning = true
			go ps.setConfigurations()
			return true
		}
		// return handlePingStart(conf, pingServiceState)
	case "stop":
		conf.Config.UtilsConf.ServicesSignal.Ping = "passive"
		_, e := conf.Write()
		if e != nil {
			panic(e)
		}
		return true
	default:
		logger.Terminal("invalid signal", "f")
	}
	return false
}

func (ps *Ping) setConfigurations() {
	pingConfig := ps.localConfig.Config.Routes
	pingInterval := ps.scrapeInterval

	urlStack := make(map[string]string)
	for _, route := range pingConfig {
		url := route.URL
		urlHash := utils.GetHash(url)
		// maintain urls uniquely
		_, ok := urlStack[urlHash]
		if !ok {
			urlStack[urlHash] = *filters.HTTPPingFilter(&url)
		}
	}

	ps.perform(urlStack, pingInterval)
}

func (ps *Ping) perform(urlStack map[string]string, pingInterval TestInterval) {
	i := 0
	config := ps.localConfig

	for {
		select {
		case <-ps.signalStop:
			break
		default:
		}

		i++
		config.Refresh()

		switch config.Config.UtilsConf.ServicesSignal.Ping {
		case "active":
			err, _ := utils.VerifyConnection()
			if !err {
				logger.Terminal("Not able to connect to externel network please check you internet connection", "p")
			} else {
				var wg sync.WaitGroup
				wg.Add(len(urlStack))
				for _, u := range urlStack {
					go ps.ping(u, 10, u, &wg)
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

		intrv := time.Duration(pingInterval.Duration)
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

func (ps *Ping) ping(urlRaw string, packets int, tsdbNameHash string, wg *sync.WaitGroup) {
	chain := ps.chain
	tsdbNameHash = utils.PathPing + "/" + "chunk_ping_" + tsdbNameHash + ".json"

	resp, err := utils.CLIPing(urlRaw, packets)
	if err != nil {
		msg := "Failure occured for " + urlRaw
		logger.Terminal(msg, "p")
		wg.Done()
		return
	}

	result := scrap.CLIPingScrap(resp)
	newBlock := *tsdb.GetNewBlock("ping", getNormalizedBlockString(*result))
	urlExists := false

	for index := range *chain {
		if (*chain)[index].Path == tsdbNameHash || ps.test {
			urlExists = true
			(*chain)[index] = (*chain)[index].Append(newBlock)
			if ps.test {
				continue
			}
			break
		}
	}

	if !urlExists {
		panic("faulty hashing! impossible to look for a hash match.")
	}

	wg.Done()
}

func getNormalizedBlockString(v utils.TypePingScrap) string {
	return fToS(v.Min) + tsdb.BlockDataSeparator + fToS(v.Avg) +
		tsdb.BlockDataSeparator + fToS(v.Max) + tsdb.BlockDataSeparator + fToS(v.Mdev)
}

func fToS(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
