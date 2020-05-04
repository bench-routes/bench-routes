package decode

import (
	"fmt"
	"strings"

	"github.com/zairza-cetb/bench-routes/src/lib/modules/jitter"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/ping"
	"github.com/zairza-cetb/bench-routes/src/metrics/system"
)

// System converts the block into Response type for easily http based
// response.
func systemDecode(block string) system.Response {
	arr := strings.Split(block, "|")
	if len(arr) != 8 {
		panic(fmt.Errorf("Invalid block segments length: Segments must be 8 in number: length: %d", len(arr)))
	}

	return system.Response{
		CPUTotalUsage: arr[0],
		Memory: system.MemoryStatsStringified{
			Total:       arr[1],
			Available:   arr[2],
			Used:        arr[3],
			UsedPercent: arr[4],
			Free:        arr[5],
		},
		Disk: system.DiskStatsStringified{
			DiskIO: arr[6],
			Cached: arr[7],
		},
	}
}

func pingDecode(block string) ping.Response {
	arr := strings.Split(block, "|")
	if len(arr) != 4 {
		panic(fmt.Errorf("Invalid block segments length: Segments must be 8 in number: length: %d", len(arr)))
	}
	return ping.Response{
		Min:  arr[0],
		Avg:  arr[1],
		Max:  arr[2],
		Mdev: arr[3],
	}
}

func jitterDecode(block string) jitter.Response {
	arr := strings.Split(block, "|")
	return jitter.Response{
		Value: arr[0],
	}
}

func floodPingDecode(block string) ping.FloodPingResponse {
	arr := strings.Split(block, "|")
	if len(arr) != 5 {
		panic(fmt.Errorf("Invalid block segments length: Segments must be 8 in number: length: %d", len(arr)))
	}
	return ping.FloodPingResponse{
		Min:   arr[0],
		Avg:   arr[1],
		Max:   arr[2],
		Mdev:  arr[3],
		Ploss: arr[4],
	}
}
