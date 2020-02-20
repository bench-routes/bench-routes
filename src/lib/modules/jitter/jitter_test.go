package jitter

import (
	"log"
	"sync"
	"testing"

	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/parser"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
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
	pathJitter            = "../../../../storage/jitter"
	ConfigURLs            []string
)

func initTest() {
	log.Printf("Testing jitter")

	// load configuration file
	Configuration.Address = ConfigurationFilePath
	Configuration = *Configuration.Load()

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

	//initiating the jitter chain
	for _, v := range ConfigURLs {
		path := pathJitter + "/" + "chunk_jitter_" + utils.PingDBNames[v] + ".json"
		inst := &tsdb.Chain{
			Path:           path,
			Chain:          []tsdb.Block{},
			LengthElements: 0,
			Size:           0,
		}
		inst.Init().Commit()
		// Initiate the chain
		utils.Jitterc = append(utils.Jitterc, inst)
	}
}

func TestHandlerJitter(t *testing.T) {

	initTest()
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, inst := range urls {
		go HandleJitter(utils.Jitterc, inst, 10, inst, &wg, true)
	}

	wg.Wait()
}
