package ping

import (
	"sync"
	"testing"

	"github.com/zairza-cetb/bench-routes/src/lib/parser"
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
	Configuration         parser.YAMLBenchRoutesType
	ConfigurationFilePath = "../../../../local-config.yml"
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
			utils.PingDBNames[r.URL] = utils.GetHash(r.URL)
		}
	}
	// forming ping chain
	for _, v := range ConfigURLs {
		path := pathPing + "/" + "chunk_ping_" + v + ".json"
		inst := &tsdb.Chain{
			Path:           path,
			Chain:          []tsdb.Block{},
			LengthElements: 0,
			Size:           0,
		}
		inst.Init().Commit()
		// Initiate the chain
		utils.Pingc = append(utils.Pingc, inst)
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
			utils.FloodPingDBNames[r.URL] = utils.GetHash(r.URL)
		}
	}
	// forming ping chain
	for _, v := range ConfigURLs {
		path := pathFloodPing + "/" + "chunk_flood_ping_" + v + ".json"
		inst := &tsdb.Chain{
			Path:           path,
			Chain:          []tsdb.Block{},
			LengthElements: 0,
			Size:           0,
		}
		inst.Init().Commit()
		// Initiate the chain
		utils.FPingc = append(utils.FPingc, inst)
	}
}

func TestHandlerPing(t *testing.T) {
	initPingTest()
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, inst := range urls {
		go HandlePing(utils.Pingc, inst, 10, "", &wg, true)
	}
	wg.Wait()
}

func TestHandlerFloodPing(t *testing.T) {
	initFloodPingTest()
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, inst := range urls {
		go HandleFloodPing(utils.FPingc, inst, 1000, "", &wg, true, Configuration.Config.Password)
	}
	wg.Wait()
}
