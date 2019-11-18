package ping

import (
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
	"log"
	"sync"
)

// HandleFloodPing is the main handler for flood ping operations
func HandleFloodPing(globalChain []*tsdb.Chain, urlRaw string, packets int, tsdbNameHash string, wg *sync.WaitGroup, isTest bool, password string) {

	tsdbNameHash = utils.PathFloodPing + "/" + "chunk_flood_ping_" + tsdbNameHash + ".json"
	resp, err := utils.CLIFloodPing(urlRaw, packets, password)
	if err != nil {
		log.Println(*resp)
		wg.Done()
		return
	}
	result := *scrap.CLIFLoodPingScrap(resp)
	block := *tsdb.GetNewBlock("flood-ping", getNormalizedBlockStringFlood(result))
	urlExists := false
	for index := range globalChain {
		if globalChain[index].Path == tsdbNameHash {
			urlExists = true
			globalChain[index] = globalChain[index].Append(block)
			globalChain[index].Save()
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
