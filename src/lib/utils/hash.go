package utils

import (
	// "crypto/sha1"
	"hash/fnv"
)

// GetHash returns an unique hash code which can be used for storing values in tsdb for long urls
func GetHash(s *string) string {
	en := fnv.New32a()
	_, err := en.Write([]byte(*s))
	if err != nil {
		panic(err)
	}

	return string(en.Sum32())
}
