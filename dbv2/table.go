package dbv2

import (
	"fmt"
	"os"
	"sync"
)

type DataTable struct {
	*os.File
	path              string
	mux               *sync.RWMutex
	offset            int16
	minWriteTimestamp uint64
	buffer            *TableBuffer
}

// OpenRWDataTable opens or create a new DataTable that can be either read and written. It creates a new DataTable
// at the specified path if a table does not exist there and returns true.
func OpenRWDataTable(path string) (dtbl *DataTable, isCreated bool, err error) {
	var file *os.File
	// Verify if file exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return dtbl, false, fmt.Errorf("create data-table: creating new data-table file: %w", err)
		}
		isCreated = true
	} else {
		file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return dtbl, isCreated, fmt.Errorf("create data-table: open existing data-table file: %w", err)
		}

	}
	dtbl = &DataTable{
		path:              path,
		File:              file,
		mux:               new(sync.RWMutex),
		offset:            3,
		minWriteTimestamp: 0,
	}
	dtbl.buffer = NewTableBuffer(dtbl, dtbl.mux, 200, 1000)
	return
}