package jitter

import (
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"strconv"
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/parser"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

// Jitter is the structure that implements the Jitter service.
type Jitter struct {
	localConfig    *parser.YAMLBenchRoutesType
	scrapeInterval TestInterval
	chain          []*tsdb.Chain
	test           bool
}

//TestInterval stores the value of the duration and the type of test
type TestInterval struct {
	OfType   string
	Duration int64
}

// New returns a Jitter type.
func New(configuration *parser.YAMLBenchRoutesType, scrapeInterval TestInterval, chain []*tsdb.Chain) *Jitter {
	return &Jitter{
		localConfig:    configuration,
		scrapeInterval: scrapeInterval,
		chain:          chain,
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

	conf := ps.localConfig
	conf.Refresh()
	pingServiceState := conf.Config.UtilsConf.ServicesSignal.Jitter

	switch signal {
	case "start":
		if pingServiceState == "passive" {
			conf.Config.UtilsConf.ServicesSignal.Jitter = "active"
			_, e := conf.Write()
			if e != nil {
				panic(e)
			}
			go ps.setConfigurations()
			return true
		}
		return true
	case "stop":
		conf.Config.UtilsConf.ServicesSignal.Jitter = "passive"
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

func (ps *Jitter) setConfigurations() {
	config := ps.localConfig.Config.Routes
	interval := ps.scrapeInterval

	urlStack := make(map[string]string)
	for _, route := range config {
		url := route.URL
		urlHash := utils.GetHash(url)
		// maintain urls uniquely
		_, ok := urlStack[urlHash]
		if !ok {
			urlStack[urlHash] = *filters.HTTPPingFilter(&url)
		}
	}

	ps.perform(urlStack, interval)
}

func (ps *Jitter) perform(urlStack map[string]string, pingInterval TestInterval) {
	i := 0
	config := ps.localConfig

	for {
		i++
		config.Refresh()

		switch config.Config.UtilsConf.ServicesSignal.Jitter {
		case "active":
			err, _ := utils.VerifyConnection()
			if !err {
				logger.Terminal("Not able to connect to externel network please check you internet connection", "p")
			} else {
				var wg sync.WaitGroup
				wg.Add(len(urlStack))
				for _, u := range urlStack {
					go ps.jitter(u, 10, u, &wg, false)
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

		intrv := time.Duration(pingInterval.Duration)
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

func (ps *Jitter) jitter(urlRaw string, packets int, tsdbNameHash string, wg *sync.WaitGroup, isTest bool) {
	chain := ps.chain
	tsdbNameHash = utils.PathJitter + "/" + "chunk_jitter_" + tsdbNameHash + ".json"

	resp, err := utils.CLIPing(urlRaw, packets)
	if err != nil {
		msg := "Failure occured for " + urlRaw
		logger.Terminal(msg, "p")
		wg.Done()
		return
	}

	result := scrap.CLIJitterScrap(resp)
	newBlock := *tsdb.GetNewBlock("jitter", fToS(result))
	urlExists := false
	for index := range chain {
		if chain[index].Path == tsdbNameHash || ps.test {
			urlExists = true
			chain[index].Append(newBlock)
			if ps.test {
				continue
			}
			break
		}
	}

	if !urlExists && !isTest {
		panic("faulty hashing! impossible to look for a hash match.")
	}

	wg.Done()
}

func fToS(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
