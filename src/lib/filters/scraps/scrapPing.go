package scraps

import (
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"strconv"
	"strings"
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
		if words[0] == "rtt" {
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
	sentence := arr[3]
	words := strings.Split(sentence, ", ")
	packetLoss := strings.Split(words[2], " ")[0]

	pingValuesString := strings.Split(arr[4], " ")[3]
	pingValues := strings.Split(pingValuesString, "/")
	a = &utils.TypeFloodPingScrap{
		Min:        strToFloat64(pingValues[0]),
		Avg:        strToFloat64(pingValues[1]),
		Max:        strToFloat64(pingValues[2]),
		Mdev:       strToFloat64(pingValues[3]),
		PacketLoss: strToFloat64(strings.TrimRight(packetLoss, "%")),
	}

	return
}
