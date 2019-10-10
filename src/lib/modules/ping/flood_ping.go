package ping

import (
	"sync"
	"time"

	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

const (
	// PathFloodPing stores the default address of storage directory of flood ping data
	PathFloodPing = "storage/flood-ping"
)

// HandleFloodPing is the main handler for flood ping operations
func HandleFloodPing(globalChain []*tsdb.ChainFloodPing, urlRaw string, packets int, tsdbNameHash string, wg *sync.WaitGroup, isTest bool, password string) {
	chnl := make(chan *string)

	tsdbNameHash = PathFloodPing + "/" + "chunk_flood_ping_" + tsdbNameHash + ".json"
	// launch a goroutine to handle ping operations
	go utils.CLIFloodPing(urlRaw, packets, chnl, password)
	resp := <-chnl
	result := *scrap.CLIFLoodPingScrap(resp)
	newBlock := createNewBlockFloodPing(result)
	urlExists := false
	for index := range globalChain {
		if globalChain[index].Path == tsdbNameHash {
			urlExists = true
			globalChain[index] = globalChain[index].AppendFloodPing(newBlock)
			globalChain[index].SaveFloodPing()
			break
		}
	}
	if !urlExists && !isTest {
		panic("faulty hashing! impossible to look for a hash match.")
	}
	wg.Done()
}

func createNewBlockFloodPing(val scrap.TypeFloodPingScrap) tsdb.BlockFloodPing {
	return tsdb.BlockFloodPing{
		Timestamp: time.Now(),
		Datapoint: tsdb.FloodPingType{
			Min:        val.Min,
			Mean:       val.Avg,
			Max:        val.Max,
			MDev:       val.Mdev,
			PacketLoss: val.PacketLoss,
		},
	}
}
