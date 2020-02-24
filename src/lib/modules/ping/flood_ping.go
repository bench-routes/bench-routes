package ping

import (
	"sync"

	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

// HandleFloodPing is the main handler for flood ping operations
func HandleFloodPing(Jitterc []*tsdb.Chain, urlRaw string, packets int, tsdbNameHash string, wg *sync.WaitGroup, isTest bool, password string) {

	tsdbNameHash = utils.PathFloodPing + "/" + "chunk_flood_ping_" + tsdbNameHash + ".json"
	resp, err := utils.CLIFloodPing(urlRaw, packets, password)
	if err != nil {
		logger.File(*resp, "p")
		wg.Done()
		return
	}
	result := *scrap.CLIFLoodPingScrap(resp)
	block := *tsdb.GetNewBlock("flood-ping", getNormalizedBlockStringFlood(result))
	urlExists := false
	for index := range Jitterc {
		if Jitterc[index].Path == tsdbNameHash {
			urlExists = true
			Jitterc[index] = Jitterc[index].Append(block)
			Jitterc[index].Commit()
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
