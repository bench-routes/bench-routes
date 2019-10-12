package jitter

import (
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
	"log"
	"sync"
	"time"
)

const (
	// PathJitter stores the defualt address of storage directory of jitter data
	PathJitter = "storage/jitter"
)

//HandleJitter handles the url and calculate the jitter for that url
func HandleJitter(globalChain []*tsdb.Chain, url string, packets int, tsdbNameHash string, wg *sync.WaitGroup, isTest bool) {
	tsdbNameHash = PathJitter + "/" + "chunk_jitter_" + tsdbNameHash + ".json"
	resp, err := utils.CLIPing(url, packets)
	if err != nil {
		log.Println(*resp)
		wg.Done()
		return
	}
	result := scrap.CLIJitterScrap(resp)
	newBlock := createNewBlock(result)
	urlExists := false
	for index := range globalChain {
		if globalChain[index].Path == tsdbNameHash {
			urlExists = true
			globalChain[index] = globalChain[index].Append(newBlock)
			globalChain[index].Save()
			break
		}
	}
	if !urlExists && !isTest {
		panic("faulty hashing! impossible to look for a hash match.")
	}
	wg.Done()
}

func createNewBlock(val float64) tsdb.Block {
	return tsdb.Block{
		Timestamp: time.Now(),
		Datapoint: float32(val),
	}
}
