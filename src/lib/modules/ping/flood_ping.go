package ping

import (
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
)

// HandleFloodPing is the main handler for flood ping operations
func HandleFloodPing(url *string, packets int) (float64, *scrap.TypePingScrap) {
	chnl := make(chan *string)

	// launch a goroutine to handle ping operations
	go utils.CLIFloodPing(url, packets, chnl)
	resp := <-chnl
	return scrap.CLIFLoodPingScrap(resp)
}
