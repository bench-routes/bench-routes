package ping

import (
	"fmt"
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
	"sync"
	"time"
)

const (
	// PathPing stores the defualt address of storage directory of ping data
	PathPing = "storage/ping"
)

// HandlePing is the main handler for ping operations
func HandlePing(globalChain []*tsdb.ChainPing, urlRaw *string, packets int, tsdbNameHash string, wg *sync.WaitGroup, isTest bool) {
	chnl := make(chan *string)

	tsdbNameHash = PathPing + "/" + "chunk_ping_" + tsdbNameHash + ".json"
	// launch a goroutine to handle ping operations
	go utils.CLIPing(urlRaw, packets, chnl)
	resp := <-chnl
	result := *scrap.CLIPingScrap(resp)
	newBlock := createNewBlock(result)
	urlExists := false
	fmt.Println(&(globalChain))
	fmt.Println(&(globalChain[0].Path))
	for index := range globalChain {
		fmt.Println("tsdbNameHash:::", tsdbNameHash, "  globalChain[index].Path:::", globalChain[index].Path)
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

func createNewBlock(val scrap.TypePingScrap) tsdb.BlockPing {
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
