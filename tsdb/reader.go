package tsdb

import (
	"errors"
	"io/ioutil"
	"sync"
)

// Reader is used to read persistent time-series samples from the storage
// as required by the user. This is go-routine safe and hence can handle multiple
// read request on the same samples as well.
type Reader struct {
	mux sync.RWMutex
}

// Open opens a reader with the default byte data stream on the specified path.
func (r *Reader) open(service, path string) (string, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	res, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.New("db not existing at the specified path: " + path)
	}
	str := string(res)

	return str, nil
}
