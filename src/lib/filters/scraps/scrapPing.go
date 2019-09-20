package scraps

import (
	"strconv"
	"strings"
)

// TypePingScrap as datatype for ping outputs
type TypePingScrap struct {
	Min, Avg, Max, Mdev float64
}

func strToFloat64(s string) float64 {
	r, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return r
}

// CLIPingScrap scraps the data points for CLIPing function
func CLIPingScrap(s *string) (a *TypePingScrap) {
	arr := strings.Split(*s, "\n")
	for _, lines := range arr {
		words := strings.Split(lines, " ")
		if words[0] == "rtt" {
			temp := strings.Split(words[len(words)-2], "/")
			a = &TypePingScrap{
				Min:  strToFloat64(temp[0]),
				Avg:  strToFloat64(temp[1]),
				Max:  strToFloat64(temp[2]),
				Mdev: strToFloat64(temp[3]),
			}
		}
	}
	return
}
