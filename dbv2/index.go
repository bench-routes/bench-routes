package dbv2

import (
	"os"
	"sync"
)

type index struct {
	mux  *sync.RWMutex
	file *os.File
	data map[uint64][]uint64
}

type indexWriter struct {
	mux *sync.RWMutex
}
