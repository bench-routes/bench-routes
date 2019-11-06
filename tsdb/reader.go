package tsdb

import (
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

	b, err := parse(path)
	if err != nil {
		return "", err
	}

	return *b, nil
}
