package ping

import (
	"fmt"
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
	"log"
	"sync"
	"time"
)

const (
	// PathPing stores the defualt address of storage directory of ping data
	PathPing = "storage/ping"
)

// HandlePing is the main handler for ping operations
func HandlePing(urlRaw *string, packets int, tsdbNameHash string, wg *sync.WaitGroup) {
	chnl := make(chan *string)

	tsdbNameHash = PathPing + "/" + "chunk_ping_" + tsdbNameHash + ".json"
	// launch a goroutine to handle ping operations
	go utils.CLIPing(urlRaw, packets, chnl)
	resp := <-chnl
	result := *scrap.CLIPingScrap(resp)
	fmt.Println(result)
	log.Printf("writing to tsdb")
	newBlock := createNewBlock(result)
	for index := range tsdb.GlobalPingChain {
		fmt.Println(tsdb.GlobalPingChain[index].Path, "  ", tsdbNameHash)
		if tsdb.GlobalPingChain[index].Path == tsdbNameHash {
			fmt.Println("match found")
			tsdb.GlobalPingChain[index] = tsdb.GlobalPingChain[index].AppendPing(newBlock)
			tsdb.GlobalPingChain[index].SavePing()
		}
	}
	// tsdb.GlobalPingChain[tsdbNameHash].SavePing()
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
