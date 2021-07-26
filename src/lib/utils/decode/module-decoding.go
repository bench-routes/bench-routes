package decode

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bench-routes/bench-routes/src/lib/modules/evaluate"
)

func pingDecode(block string) evaluate.Ping {
	if block == "" {
		return evaluate.Ping{}
	}
	arr := strings.Split(block, "|")
	if len(arr) != 3 {
		panic(fmt.Errorf("invalid block segments length: Segments must be 3 in number: length: %d", len(arr)))
	}

	min, err := strconv.Atoi(arr[0])
	if err != nil {
		panic(fmt.Errorf("error parsing datapoint of type Ping Min"))
	}
	mean, err := strconv.Atoi(arr[1])
	if err != nil {
		panic(fmt.Errorf("error parsing datapoint of type Ping Mean"))
	}
	max, err := strconv.Atoi(arr[2])
	if err != nil {
		panic(fmt.Errorf("error parsing datapoint of type Ping Max"))
	}
	return evaluate.Ping{
		Min:  time.Millisecond * time.Duration(min/1000),
		Mean: time.Millisecond * time.Duration(mean/1000),
		Max:  time.Millisecond * time.Duration(max/1000),
	}
}

func jitterDecode(block string) evaluate.Jitter {
	if block == "" {
		return evaluate.Jitter{}
	}
	arr := strings.Split(block, "|")
	if len(arr) != 1 {
		panic(fmt.Errorf("invalid block segments length: Segments must be 1 in number: length: %d", len(arr)))
	}

	jitter, err := strconv.Atoi(arr[0])
	if err != nil {
		panic(fmt.Errorf("error parsing datapoint of type Jitter"))
	}

	return evaluate.Jitter{
		Value: time.Millisecond * time.Duration(jitter/1000),
	}
}

func monitorDecode(block string) evaluate.Response {
	if block == "" {
		return evaluate.Response{}
	}
	arr := strings.Split(block, "|")
	if len(arr) != 4 {
		panic(fmt.Errorf("invalid block segments length: Segments must be 4 in number: length: %d", len(arr)))
	}

	delay, err := strconv.Atoi(arr[0])
	if err != nil {
		panic(fmt.Errorf("error parsing Delay datapoint"))
	}
	length, err := strconv.Atoi(arr[1])
	if err != nil {
		panic(fmt.Errorf("error parsing Length datapoint"))
	}
	size, err := strconv.Atoi(arr[2])
	if err != nil {
		panic(fmt.Errorf("error parsing Size datapoint"))
	}
	status, err := strconv.Atoi(arr[3])
	if err != nil {
		panic(fmt.Errorf("error parsing Status datapoint"))
	}
	return evaluate.Response{
		Delay:  time.Millisecond * time.Duration(delay/1000),
		Length: length,
		Size:   size,
		Status: status,
	}
}
