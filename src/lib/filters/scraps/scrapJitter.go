package scraps

import (
	"math"
	"strings"
)

// CLIJitterScrap scraps the data points for HandleJitter function
func CLIJitterScrap(s *string) (jitter *float64) {
	arr := strings.Split(*s, "\n")
	var time []float64
	for _, lines := range arr {
		words := strings.Split(lines, " ")
		searchWord := "time="
		if len(words) >= 8 {
			if strings.Contains(words[7], searchWord) {
				found := strings.Split(words[7], "")
				timeConsumed := strToFloat64(strings.Join(found[5:], ""))
				time = append(time, timeConsumed)
			}
		}
	}
	jitter = calculateJitter(time)
	return
}

func calculateJitter(timeArr []float64) (sum *float64) {
	x := 0.0
	sum = &x
	for i := 1; i < len(timeArr); i++ {
		*sum += math.Abs(timeArr[i] - timeArr[i-1])
	}
	*sum = *sum / float64(len(timeArr)-1)
	return
}
