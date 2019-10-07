package ping

import (
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
	"log"
	"reflect"
	"sync"
	"testing"

	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
)

var (
	urls = []string{
		"google.co.in",
		"facebook.com",
		"yahoo.com",
		"youtube.com",
	}
	Configuration         utils.YAMLBenchRoutesType
	ConfigurationFilePath = "../../../../local-config.yml"
	pathPing              = "../../../../storage/ping"
	ConfigURLs            []string
)

func initTest() {
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

func TestHandlerPing(t *testing.T) {
	initTest()
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, inst := range urls {
		go HandlePing(tsdb.GlobalPingChain, &inst, 10, "", &wg, true)
	}
	wg.Wait()
}

func TestHandlerFloodPing(t *testing.T) {
	for _, inst := range urls {
		a, b := HandleFloodPing(&inst, 1000)
		if reflect.TypeOf(a).String() != "float64" {
			t.Errorf("invalid packet loss type, not a float64 value")
		} else if a > 100 || a < 0 {
			t.Errorf("invalid packet response")
		} else {
			log.Println(a)
		}
		if *b == (scrap.TypePingScrap{}) {
			t.Errorf("invalid ping response")
		} else {
			log.Println(*b)
		}
	}
}
