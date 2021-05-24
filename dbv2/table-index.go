package dbv2

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

const indexFileExtension = "index"

// tableIndex is an index between a time-series id, corresponding to its list of indices.
// This is represented in uncompressed form on disk as:
// 1: [1, 7, 13, 19, 29, 30]
// 2: [3, 8, 15, 23, 37, 40]
// this goes on till the number of time-series in the system for that index.
type tableIndex struct {
	mux   sync.RWMutex
	path  string
	index map[uint64][]uint64
}

func NewTableIndex(path string) *tableIndex {
	return &tableIndex{
		path:  fmt.Sprintf("%s.%s", path, indexFileExtension),
		index: make(map[uint64][]uint64),
	}
}

// Add adds the gap between the byte position of previous value of supplied series_id and the current insertion.
func (ti *tableIndex) Add(series_id uint64, currentByteVal uint64) error {
	ti.mux.Lock()
	defer ti.mux.Unlock()
	if _, ok := ti.index[series_id]; !ok {
		ti.index[series_id] = []uint64{currentByteVal}
		return nil
	}
	l := len(ti.index[series_id])
	if currentByteVal <= ti.index[series_id][l-1] {
		// Index of ending position of insertion into table is always increasing in order. If this is not found,
		// we should error the insert as this would be a faulty insert.
		return fmt.Errorf(
			"inserting endAtPos should always be greater than the previous endAtPos for that series: expected pos >%d, received pos %d",
			ti.index[series_id][l-1],
			currentByteVal,
		)
	}
	ti.index[series_id] = append(ti.index[series_id], currentByteVal)
	return nil
}

// Get returns the ending positions of the entries. If not found, it returns an empty array.
func (ti *tableIndex) Get(series_id uint64) (positionGapIndices []uint64) {
	ti.mux.RLock()
	defer ti.mux.RUnlock()
	if positionGapIndices, ok := ti.index[series_id]; ok {
		return positionGapIndices
	}
	return []uint64{}
}

// Flush flushes the in-memory index to disk.
func (ti *tableIndex) Flush() error {
	f, err := os.Create(ti.path)
	if err != nil {
		return fmt.Errorf("creating index-file: %w", err)
	}
	defer f.Close()
	for k, v := range ti.index {
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
func (ti *tableIndex) Load() error {
	// todo: needs a test
	f, err := os.Open(ti.path)
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
		indexSep := bytes.Index(b, []byte{colon})
		idBytes := b[:indexSep]
		indicesBytes := b[indexSep+1:]
		var (
			id      uint64
			indices []uint64
		)
		if err := json.Unmarshal(idBytes, &id); err != nil {
			return fmt.Errorf("decoding id: %w", err)
		}
		if err := json.Unmarshal(indicesBytes, &indices); err != nil {
			return fmt.Errorf("decoding indices: %w", err)
		}
		ti.index[id] = indices
	}
	return nil
}
