package dbv2

import (
	"fmt"
	"sync"
)

type tableIndex struct {
	mux   sync.RWMutex
	index map[uint64][]uint64
}

func NewTableIndex() *tableIndex {
	return &tableIndex{
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
