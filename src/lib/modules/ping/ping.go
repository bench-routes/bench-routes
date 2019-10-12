package ping

import (
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
	"log"
	"sync"
	"time"
)

// const (
// 	// PathPing stores the defualt address of storage directory of ping data
// 	PathPing = "storage/ping"
// )

// HandlePing is the main handler for ping operations
func HandlePing(globalChain []*tsdb.ChainPing, urlRaw string, packets int, tsdbNameHash string, wg *sync.WaitGroup, isTest bool) {
	tsdbNameHash = utils.PathPing + "/" + "chunk_ping_" + tsdbNameHash + ".json"
	// launch a goroutine to handle ping operations
	resp, err := utils.CLIPing(urlRaw, packets)
	if err != nil {
		log.Println(*resp)
		wg.Done()
		return
	}
	result := *scrap.CLIPingScrap(resp)
	newBlock := createNewBlock(result)
	urlExists := false
	for index := range globalChain {
		if globalChain[index].Path == tsdbNameHash {
			urlExists = true
			globalChain[index] = globalChain[index].AppendPing(newBlock)
			globalChain[index].SavePing()
			break
		}
	}
	if !urlExists && !isTest {
		panic("faulty hashing! impossible to look for a hash match.")
	}
	wg.Done()
}

func createNewBlock(val utils.TypePingScrap) tsdb.BlockPing {
	return tsdb.BlockPing{
		Timestamp: time.Now(),
		Datapoint: tsdb.PingType{
			Min:  val.Min,
			Mean: val.Avg,
			Max:  val.Max,
			MDev: val.Mdev,
		},
	}
}
