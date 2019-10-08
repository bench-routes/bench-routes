package ping

import (
	"sync"
	"testing"

	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"

	"github.com/zairza-cetb/bench-routes/src/lib/filters"
)

var (
	urls = []string{
		"google.co.in",
		"facebook.com",
		"yahoo.com",
		"youtube.com",
	}
	Configuration         utils.YAMLBenchRoutesType
	ConfigurationFilePath = "../../../../storage/local-config.yml"
	pathPing              = "../../../../storage/ping"
	pathFloodPing         = "../../../../storage/flood-ping"
)

func initPingTest() {
	Configuration.Address = ConfigurationFilePath
	Configuration = *Configuration.Load()

	var ConfigURLs []string

	// Load and build TSDB chain
	// searching for unique URLs
	for _, r := range Configuration.Config.Routes {
		found := false
		for _, i := range ConfigURLs {
			if i == r.URL {
				found = true
				break
			}
		}
		if !found {
			filters.HTTPPingFilter(&r.URL)
			ConfigURLs = append(ConfigURLs, r.URL)
			tsdb.PingDBNames[r.URL] = utils.GetHash(r.URL)
		}
	}
	// forming ping chain
	for i, v := range ConfigURLs {
		path := pathPing + "/" + "chunk_ping_" + tsdb.PingDBNames[v] + ".json"
		inst := &tsdb.ChainPing{
			Path:           path,
			Chain:          []tsdb.BlockPing{},
			LengthElements: 0,
			Size:           0,
		}
		// Initiate the chain
		tsdb.GlobalPingChain = append(tsdb.GlobalPingChain, inst)
		tsdb.GlobalPingChain[i] = tsdb.GlobalPingChain[i].InitPing()
		tsdb.GlobalPingChain[i].SavePing()
	}
}

func initFloodPingTest() {
	// load configuration file
	Configuration.Address = ConfigurationFilePath
	Configuration = *Configuration.Load()

	var ConfigURLs []string

	// Load and build TSDB chain
	// searching for unique URLs
	for _, r := range Configuration.Config.Routes {
		found := false
		for _, i := range ConfigURLs {
			if i == r.URL {
				found = true
				break
			}
		}
		if !found {
			filters.HTTPPingFilter(&r.URL)
			ConfigURLs = append(ConfigURLs, r.URL)
			tsdb.FloodPingDBNames[r.URL] = utils.GetHash(r.URL)
		}
	}
	// forming ping chain
	for i, v := range ConfigURLs {
		path := pathFloodPing + "/" + "chunk_flood_ping_" + v + ".json"
		inst := &tsdb.ChainFloodPing{
			Path:           path,
			Chain:          []tsdb.BlockFloodPing{},
			LengthElements: 0,
			Size:           0,
		}
		// Initiate the chain
		tsdb.GlobalFloodPingChain = append(tsdb.GlobalFloodPingChain, inst)
		tsdb.GlobalFloodPingChain[i] = tsdb.GlobalFloodPingChain[i].InitFloodPing()
		tsdb.GlobalFloodPingChain[i].SaveFloodPing()
	}
}

func TestHandlerPing(t *testing.T) {
	initPingTest()
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, inst := range urls {
		go HandlePing(tsdb.GlobalPingChain, &inst, 10, "", &wg, true)
	}
	wg.Wait()
}

func TestHandlerFloodPing(t *testing.T) {
	initFloodPingTest()
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, inst := range urls {
		go HandleFloodPing(tsdb.GlobalFloodPingChain, &inst, 1000, "", &wg, true, Configuration.Config.Password)
	}
	wg.Wait()
}
