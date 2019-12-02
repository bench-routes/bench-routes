package utils

import (
	"fmt"
	"hash/fnv"
	"github.com/zairza-cetb/bench-routes/tsdb"
	"strings"
	"strconv"
	"errors"
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
			Min: sTof(arr[0]),
			Mean: sTof(arr[1]),
			Max: sTof(arr[2]),
			MDev: sTof(arr[3]),
		}
	case "jitter":
		return sTof(b.GetDatapointEnc())
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
