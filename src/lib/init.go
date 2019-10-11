package lib

import (
	"log"

	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

func init() {

	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.Printf("initializing bench-routes ...")

	// load configuration file
	Configuration.Address = utils.ConfigurationFilePath
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
			tsdb.FloodPingDBNames[r.URL] = utils.GetHash(r.URL)
		}
	}
	// forming ping chain
	for i, v := range ConfigURLs {
		path := utils.PathPing + "/" + "chunk_ping_" + v + ".json"
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

	// forming ping chain
	for i, v := range ConfigURLs {
		path := utils.PathFloodPing + "/" + "chunk_flood_ping_" + v + ".json"
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

	for i, v := range ConfigURLs {
		path := utils.PathJitter + "/" + "chunk_jitter_" + v + ".json"
		inst := &tsdb.Chain{
			Path:           path,
			Chain:          []tsdb.Block{},
			LengthElements: 0,
			Size:           0,
		}
		// Initiate the chain
		tsdb.GlobalChain = append(tsdb.GlobalChain, inst)
		tsdb.GlobalChain[i] = tsdb.GlobalChain[i].Init()
		tsdb.GlobalChain[i].Save()
	}

	// forming req-res-delay chain
	for i, route := range Configuration.Config.Routes {
		path := utils.PathReqResDelayMonitoring + "/" + "chunk_req_res_" + filters.RouteDestroyer(route.URL)
		// Create sample chains to init in each TSDB file
		sampleResponseDelay := &tsdb.Chain{
			Path:           path + "_delay.json",
			Chain:          []tsdb.Block{},
			LengthElements: 0,
			Size:           0,
		}
		sampleResponseLength := &tsdb.Chain{
			Path:           path + "_length.json",
			Chain:          []tsdb.Block{},
			LengthElements: 0,
			Size:           0,
		}
		sampleResponseStatusCode := &tsdb.Chain{
			Path:           path + "_status.json",
			Chain:          []tsdb.Block{},
			LengthElements: 0,
			Size:           0,
		}
		tsdb.GlobalResponseDelay = append(tsdb.GlobalResponseDelay, sampleResponseDelay)
		tsdb.GlobalResponseLength = append(tsdb.GlobalResponseLength, sampleResponseLength)
		tsdb.GlobalResponseStatusCode = append(tsdb.GlobalResponseStatusCode, sampleResponseStatusCode)

		// Initiate all chains
		tsdb.GlobalResponseDelay[i] = tsdb.GlobalResponseDelay[i].Init()
		tsdb.GlobalResponseLength[i] = tsdb.GlobalResponseLength[i].Init()
		tsdb.GlobalResponseStatusCode[i] = tsdb.GlobalResponseStatusCode[i].Init()

		// Commit all chains to the TSDB
		tsdb.GlobalResponseDelay[i].Save()
		tsdb.GlobalResponseLength[i].Save()
		tsdb.GlobalResponseStatusCode[i].Save()
	}

	// keep the below line to the end of file so that we ensure that we give a confirmation message only when all the
	// required resources for the application is up and healthy
	log.Println("Bench-routes is up and running")
}
