package dbv2

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

const seriesFileExtension = "series"

var ErrSeriesRelNotFound = fmt.Errorf("series relation not found")

// seriesRelation is a relationship between time-series and id. This id is used in tableIndex. id is a monotonically
// increasing value based on new time-series seen.
type seriesRelation struct {
	mux        sync.RWMutex
	path       string
	curr       uint64
	newIndices bool
	index      map[string]uint64
}

func NewSeriesRelation(path string) *seriesRelation {
	return &seriesRelation{
		path:  fmt.Sprintf("%s.%s", path, seriesFileExtension),
		index: make(map[string]uint64),
	}
}

// Add adds the gap between the byte position of previous value of supplied series_id and the current insertion.
// it returns the timeseries id which.
func (sr *seriesRelation) Add(timeseries string) uint64 {
	sr.mux.Lock()
	defer sr.mux.Unlock()
	if !sr.newIndices {
		sr.newIndices = true
	}
	if id, ok := sr.index[timeseries]; ok {
		return id
	}
	atomic.AddUint64(&sr.curr, 1)
	sr.index[timeseries] = sr.curr
	return sr.curr
}

// Get returns the ending positions of the entries. If not found, it returns an error.
func (sr *seriesRelation) GetID(timeseries string) (uint64, error) {
	sr.mux.RLock()
	defer sr.mux.RUnlock()
	if id, ok := sr.index[timeseries]; ok {
		return id, nil
	}
	return 0, ErrSeriesRelNotFound
}

// Flush flushes the in-memory index to disk.
func (sr *seriesRelation) Flush() error {
	// todo: more performance efficient if using littleendian based byte conversion.
	sr.mux.Lock()
	defer sr.mux.Unlock()
	if !sr.newIndices {
		// Do not flush if no new indices.
		return nil
	}
	f, err := os.Create(sr.path)
	if err != nil {
		return fmt.Errorf("creating index-file: %w", err)
	}
	defer f.Close()
	for k, v := range sr.index {
		bslice, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("flush table-index: %w", err)
		}

		// Writing.
		marshalled, err := json.Marshal(k)
		if err != nil {
			return fmt.Errorf("table-index: marshalling key: %w", err)
		}
		net := marshalled
		net = append(net, valueSeparator)
		net = append(net, bslice...)
		net = append(net, newLineSymbol)
		if _, err := f.Write(net); err != nil {
			return fmt.Errorf("writing index-file: %w", err)
		}
	}
	return nil
}

// Load loads the index from the file. This is helpful in case of restart.
func (sr *seriesRelation) Load() error {
	// todo: needs a test
	f, err := os.Open(sr.path)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	// todo: replace r => NewReader with a NewScanner.
	r := bufio.NewReader(f)
	for {
		b, err := r.ReadBytes(newLineSymbol)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("load: %w", err)
		}
		indexSep := bytes.Index(b, []byte{newLineSymbol})
		idBytes := b[:indexSep]
		indicesBytes := b[indexSep+1:]
		var (
			timeseries string
			id         uint64
		)
		if err := json.Unmarshal(indicesBytes, &timeseries); err != nil {
			return fmt.Errorf("decoding series-relation: %w", err)
		}
		if err := json.Unmarshal(idBytes, &id); err != nil {
			return fmt.Errorf("decoding series-relation: %w", err)
		}
		sr.index[timeseries] = id
	}
	return nil
}
