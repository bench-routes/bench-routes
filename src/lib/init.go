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

	initPing()
	initFloodPing()

	// keep the below line to the end of file so that we ensure that we give a confirmation message only when all the
	// required resources for the application is up and healthy
	log.Println("Bench-routes is up and running")
}

func initPing() {
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
			tsdb.PingDBNames[r.URL] = utils.GetHash(r.URL)
		}
	}
	// forming ping chain
	for i, v := range ConfigURLs {
		path := PathPing + "/" + "chunk_ping_" + v + ".json"
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

func initFloodPing() {
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
		path := PathFloodPing + "/" + "chunk_flood_ping_" + v + ".json"
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
