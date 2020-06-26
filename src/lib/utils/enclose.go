package utils

import (
	"errors"
	"fmt"
	"github.com/zairza-cetb/bench-routes/tsdb"
	"hash/fnv"
	"strconv"
	"strings"
)

// GetHash returns an unique hash code which can be used for storing values in tsdb for long urls
func GetHash(s string) string {
	en := fnv.New32a()
	_, err := en.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	sstring := fmt.Sprint(en.Sum32())
	return sstring
}

// Decode returns the decoded persistent block after parsing from the stringified form.
func Decode(b tsdb.Block) interface{} {
	switch b.Type {
	case "ping":
		arr := strings.Split(b.GetDatapointEnc(), tsdb.BlockDataSeparator)
		return Ping{
			Min:  sTof(arr[0]),
			Mean: sTof(arr[1]),
			Max:  sTof(arr[2]),
			MDev: sTof(arr[3]),
		}
	case "jitter":
		return sTof(b.GetDatapointEnc())
	case "flood-ping":
		arr := strings.Split(b.GetDatapointEnc(), tsdb.BlockDataSeparator)
		return FloodPing{
			Min:        sTof(arr[0]),
			Mean:       sTof(arr[1]),
			Max:        sTof(arr[2]),
			MDev:       sTof(arr[3]),
			PacketLoss: sTof(arr[4]),
		}
	case "req-res":
		arr := strings.Split(b.GetDatapointEnc(), tsdb.BlockDataSeparator)
		return Response{
			Delay:         sTof(arr[0]),
			ResLength:     sToI(arr[1]),
			ResStatusCode: sToI(arr[2]),
		}
	}
	return errors.New("type does not match")
}

func sTof(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func sToI(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
