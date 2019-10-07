package controllers

import (
	"context"

	core "github.com/zairza-cetb/bench-routes/src/lib"
)

// PingController controllers the ping requests and transfers to the respective handler
func PingController(ctx context.Context, sig string) bool {
	return core.HandlerPingGeneral(ctx, sig)
}
