package ping

import (
	"fmt"
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"sync"
)

// HandlePing is the main handler for ping operations
func HandlePing(urlRaw *string, packets int, wg *sync.WaitGroup) {
	chnl := make(chan *string)

	// launch a goroutine to handle ping operations
	go utils.CLIPing(urlRaw, packets, chnl)
	resp := <-chnl
	wg.Done()
	fmt.Println(scrap.CLIPingScrap(resp))
}
