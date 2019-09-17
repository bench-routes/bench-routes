package jitter

import (
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
)

//HandleJitter handles the url and calculate the jitter for that url
func HandleJitter(url *string, packets int) (jitter *float64) {

	chnl := make(chan *string)
	go utils.CLIPing(url, packets, chnl)
	resp := <- chnl
	return scrap.CLIJitterScrap(resp)
}