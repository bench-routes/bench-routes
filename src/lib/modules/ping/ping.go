package ping

import (
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
	"log"
	"sync"
	"strconv"
)

// HandlePing is the main handler for ping operations
func HandlePing(chain []*tsdb.Chain, urlRaw string, packets int, tsdbNameHash string, wg *sync.WaitGroup, isTest bool) {
	tsdbNameHash = utils.PathPing + "/" + "chunk_ping_" + tsdbNameHash + ".json"

	resp, err := utils.CLIPing(urlRaw, packets)
	if err != nil {
		log.Println(*resp)
		wg.Done()
		return
	}
	result := *scrap.CLIPingScrap(resp)
	newBlock := *tsdb.GetNewBlock("ping", getNormalizedBlockString(result))
	urlExists := false
	for index := range chain {
		if chain[index].Path == tsdbNameHash {
			urlExists = true
			chain[index] = chain[index].Append(newBlock)
			chain[index].Save()
			break
		}
	}
	if !urlExists && !isTest {
		panic("faulty hashing! impossible to look for a hash match.")
	}
	wg.Done()
}

func getNormalizedBlockString(v utils.TypePingScrap) string {
	return fToS(v.Min) + tsdb.BlockDataSeparator + fToS(v.Avg) +
		tsdb.BlockDataSeparator + fToS(v.Max) + tsdb.BlockDataSeparator + fToS(v.Mdev)
}

func fToS(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
