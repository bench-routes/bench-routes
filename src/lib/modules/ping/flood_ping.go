package ping

import (
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/parser"
	"sync"
	"time"

	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

// FloodPing is the structure that implements the Ping service.
type FloodPing struct {
	localConfig    *parser.YAMLBenchRoutesType
	scrapeInterval TestInterval
	chain          []*tsdb.Chain
	password       string
	test           bool
}

// Newf returns a Flood Ping type.
func Newf(configuration *parser.YAMLBenchRoutesType, scrapeInterval TestInterval, chain []*tsdb.Chain, password string) *FloodPing {
	return &FloodPing{
		localConfig:    configuration,
		scrapeInterval: scrapeInterval,
		chain:          chain,
		password:       password,
		test:           false,
	}
}

// Iterate iterates over the local-configuration file to keep state
// of the ping service in sync with the local configuration.
// It is responsible for stopping the service without damaging the currently
// calculated samples.
func (ps *FloodPing) Iteratef(signal string, isTest bool) bool {
	if isTest {
		ps.test = true
	}
	conf := ps.localConfig
	conf.Refresh()
	pingServiceState := conf.Config.UtilsConf.ServicesSignal.FloodPing

	switch signal {
	case "start":
		if pingServiceState == "passive" {

			conf.Config.UtilsConf.ServicesSignal.FloodPing = "active"
			_, e := conf.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				ps.setConfigurations()
			}()
			return true
		}
		// return handlePingStart(conf, pingServiceState)
	case "stop":
		conf.Config.UtilsConf.ServicesSignal.FloodPing = "passive"
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

func (ps *FloodPing) setConfigurations() {
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

func (ps *FloodPing) perform(urlStack map[string]string, pingInterval TestInterval) {
	i := 0
	config := ps.localConfig

	for {
		i++
		config.Refresh()

		switch config.Config.UtilsConf.ServicesSignal.FloodPing {
		case "active":
			err, _ := utils.VerifyConnection()
			if !err {
				logger.Terminal("Not able to connect to externel network please check you internet connection", "p")
			} else {
				var wg sync.WaitGroup
				wg.Add(len(urlStack))
				for _, u := range urlStack {
					go ps.ping(u, 10, u, &wg, false)
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

func (ps *FloodPing) ping(urlRaw string, packets int, tsdbNameHash string, wg *sync.WaitGroup, isTest bool) {
	tsdbNameHash = utils.PathFloodPing + "/" + "chunk_flood_ping_" + tsdbNameHash + ".json"
	password := ps.password

	resp, err := utils.CLIFloodPing(urlRaw, packets, password)
	if err != nil {
		logger.File(*resp, "p")
		wg.Done()
		return
	}

	result := *scrap.CLIFLoodPingScrap(resp)
	newBlock := *tsdb.GetNewBlock("flood-ping", getNormalizedBlockStringFlood(result))
	urlExists := false

	c := ps.chain
	for index := range c {
		if c[index].Path == tsdbNameHash || ps.test {
			urlExists = true
			c[index] = c[index].Append(newBlock)
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

func getNormalizedBlockStringFlood(v utils.TypeFloodPingScrap) string {
	return fToS(v.Min) + tsdb.BlockDataSeparator + fToS(v.Avg) + tsdb.BlockDataSeparator +
		fToS(v.Max) + tsdb.BlockDataSeparator + fToS(v.Mdev) + tsdb.BlockDataSeparator + fToS(v.PacketLoss)
}

//// HandleFloodPing is the main handler for flood ping operations
//func HandleFloodPing(Jitterc []*tsdb.Chain, urlRaw string, packets int, tsdbNameHash string, wg *sync.WaitGroup, isTest bool, password string) {
//
//	tsdbNameHash = utils.PathFloodPing + "/" + "chunk_flood_ping_" + tsdbNameHash + ".json"
//	resp, err := utils.CLIFloodPing(urlRaw, packets, password)
//	if err != nil {
//		logger.File(*resp, "p")
//		wg.Done()
//		return
//	}
//	result := *scrap.CLIFLoodPingScrap(resp)
//	block := *tsdb.GetNewBlock("flood-ping", getNormalizedBlockStringFlood(result))
//	urlExists := false
//	for index := range Jitterc {
//		if Jitterc[index].Path == tsdbNameHash {
//			urlExists = true
//			Jitterc[index] = Jitterc[index].Append(block)
//			Jitterc[index]
//			break
//		}
//	}
//	if !urlExists && !isTest {
//		panic("faulty hashing! impossible to look for a hash match.")
//	}
//	wg.Done()
//}
