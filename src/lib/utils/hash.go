package utils

import (
	"fmt"
	"hash/fnv"
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
