package jitter

import (
	"strconv"
	"sync"

	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

const (
	// PathJitter stores the defualt address of storage directory of jitter data
	PathJitter = "storage/jitter"
)

//HandleJitter handles the url and calculate the jitter for that url
func HandleJitter(Jitterc []*tsdb.Chain, url string, packets int, tsdbNameHash string, wg *sync.WaitGroup, isTest bool) {
	tsdbNameHash = PathJitter + "/" + "chunk_jitter_" + tsdbNameHash + ".json"
	resp, err := utils.CLIPing(url, packets)
	if err != nil {
		msg := "Failure occured for" + url
		logger.Terminal(msg, "p")
		wg.Done()
		return
	}
	result := scrap.CLIJitterScrap(resp)
	newBlock := *tsdb.GetNewBlock("jitter", fToS(result))
	urlExists := false
	for index := range Jitterc {
		if Jitterc[index].Path == tsdbNameHash {
			urlExists = true
			Jitterc[index] = Jitterc[index].Append(newBlock)
			Jitterc[index].Commit()
			break
		}
	}
	if !urlExists && !isTest {
		panic("faulty hashing! impossible to look for a hash match.")
	}
	wg.Done()
}

func fToS(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
