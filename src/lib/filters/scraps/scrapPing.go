package scraps

import (
	"strconv"
	"strings"

	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
)

func strToFloat64(s string) float64 {
	r, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return r
}

// CLIPingScrap scraps the data points for CLIPing function
func CLIPingScrap(s *string) (a *utils.TypePingScrap) {
	arr := strings.Split(*s, "\n")
	l := len(arr)
	if l > 2 {
		words := strings.Split(arr[l-2], " ")
		if words[0] == "rtt" || words[0] == "round-trip" {
			temp := strings.Split(words[3], "/")
			a = &utils.TypePingScrap{
				Min:  strToFloat64(temp[0]),
				Avg:  strToFloat64(temp[1]),
				Max:  strToFloat64(temp[2]),
				Mdev: strToFloat64(temp[3]),
			}
		}
	}
	return
}

// CLIFLoodPingScrap returns packet loss
func CLIFLoodPingScrap(s *string) (a *utils.TypeFloodPingScrap) {
	arr := strings.Split(*s, "\n")
	var (
		// Packetloss line is of the form:
		// "10 packets transmitted, 2 packets received, 80.0% packet loss"
		packetLossLine string
		// PingValues line is of the form:
		// round-trip min/avg/max/stddev = 91.299/96.281/101.264/4.983 ms (Mac OS/Unix like systems)
		// rtt min/avg/max/stddev = 91.299/96.281/101.264/4.983 ms (Linux)
		pingValuesLine string
	)
	for _, line := range arr {
		if strings.Contains(line, "packet loss") {
			packetLossLine = line
		} else if strings.Contains(line, "round-trip") {
			pingValuesLine = line
		} else if strings.Contains(line, "rtt") {
			pingValuesLine = line
		}
	}
	// This error handler is for an optional security
	// if we could not find suitable parameters in the
	// ping command's output.
	if packetLossLine == "" || pingValuesLine == "" {
		logger.Terminal("p", "couldn't find a suitable response which contains packet loss")
	}
	words := strings.Split(packetLossLine, ", ")
	packetLoss := strings.Split(words[2], " ")[0]

	pingValuesString := strings.Split((strings.Split(pingValuesLine, "= ")[1]), " ms")[0]
	pingValues := strings.Split(pingValuesString, "/")
	a = &utils.TypeFloodPingScrap{
		Min:        strToFloat64(pingValues[0]),
		Avg:        strToFloat64(pingValues[1]),
		Max:        strToFloat64(pingValues[2]),
		Mdev:       strToFloat64(strings.TrimRight(pingValues[3], " ms")),
		PacketLoss: strToFloat64(strings.TrimRight(packetLoss, "%")),
	}

	return
}
