package decode

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zairza-cetb/bench-routes/src/lib/modules/jitter"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/ping"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/src/metrics/journal"
	"github.com/zairza-cetb/bench-routes/src/metrics/system"
)

// systemDecode converts the block into Response type for easy http based JSON response.
func systemDecode(block string) system.Response {
	arr := strings.Split(block, "|")
	if len(arr) != 14 {
		panic(fmt.Errorf("Invalid block segments length: Segments must be 14 in number: length: %d", len(arr)))
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
		Network: system.NetworkStats{
			PtcpIncoming: convertToInt(arr[8]),
			PtcpOutgoing: convertToInt(arr[9]),
			StcpIncoming: convertToInt(arr[10]),
			StcpOutgoing: convertToInt(arr[11]),
			PudpIncoming: convertToInt(arr[12]),
			PudpOutgoing: convertToInt(arr[13]),
		},
	}
}

func pingDecode(block string) ping.Response {
	arr := strings.Split(block, "|")
	if len(arr) != 4 {
		panic(fmt.Errorf("Invalid block segments length: Segments must be 4 in number: length: %d", len(arr)))
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
		panic(fmt.Errorf("Invalid block segments length: Segments must be 5 in number: length: %d", len(arr)))
	}
	return ping.FloodPingResponse{
		Min:   arr[0],
		Avg:   arr[1],
		Max:   arr[2],
		Mdev:  arr[3],
		Ploss: arr[4],
	}
}

func monitorDecode(block string) utils.Response {
	arr := strings.Split(block, "|")
	if len(arr) != 3 {
		panic(fmt.Errorf("Invalid block segments length: Segments must be 3 in number: length: %d", len(arr)))
	}
	return utils.Response{
		Delay:         convertToFloat64(arr[0]),
		ResLength:     convertToInt(arr[1]),
		ResStatusCode: convertToInt(arr[2]),
	}
}

func journalDecode(block string) journal.Points {
	arr := strings.Split(block, "|")
	if len(arr) != 6 {
		panic(fmt.Errorf("Invalid block segments length: Segments must be 6 in number: length: %d", len(arr)))
	}
	return journal.Decode(arr)
}

func convertToInt(s string) int {
	// Old unit tests in tsdb/querier have only 8 values and rest as empty strings.
	// In order to make that working and as a general response, we can make empty
	// strings return a zero value.
	if s == "" {
		return 0
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return v
}

func convertToFloat64(s string) float64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return v
}
