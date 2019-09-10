package ping

import (
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
)

// HandlePing is the main handler for ping operations
func HandlePing(url *string, packets int) *scrap.TypePingScrap {
	chnl := make(chan *string)

	// launch a goroutine to handle ping operations
	go utils.CLIPing(url, packets, chnl)
	resp := <- chnl
	return scrap.CLIPingScrap(resp)
}
