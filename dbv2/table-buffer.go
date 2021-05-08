package dbv2

import (
	"fmt"
	"sort"
	"sync"
)

// TableBuffer is a in-memory table that is in raw format and is yet to be commited.
type TableBuffer struct {
	bufferMux     sync.RWMutex
	data          []writeData
	cap           uint64
	idata          uint64
	tableCopy     *DataTable
	writer        *tableWriter
}

const defaultBufferSize = int(1e4)

// NewTableBuffer returns a new TableBuffer.
func NewTableBuffer(dataTable *DataTable, mux *sync.RWMutex, cap uint64, ioBufferSize int) *TableBuffer {
	return &TableBuffer{
		cap:           cap,
		idata: 0,
		data: make([]writeData, cap+10000),
		tableCopy:     dataTable,
		writer:        newTableWriter(dataTable, mux, ioBufferSize),
	}
}

type writeData struct {
	timestamp uint64
	id        uint64
	valueSet  string
}

func (tbuf *TableBuffer) Write(timestamp, id uint64, valueSet string) error {
	tbuf.bufferMux.Lock()
	defer tbuf.bufferMux.Unlock()
	row := writeData{id: id, valueSet: valueSet, timestamp: timestamp}
	tbuf.data[tbuf.idata] = row
	tbuf.idata++
	if tbuf.idata > tbuf.cap {
		if err := tbuf.flushToIOBuffer(true); err != nil {
			return fmt.Errorf("table-buffer write: %w", err)
		}
	}
	return nil
}

// flushToIOBuffer is not concurrent safe. The caller is expected to use a lock.
// TODO: this can be problematic in situations when the API takes more than flushDeadline. This can be resolved by
// increasing the flushDuration to an hour and maintain the state of table buffer in a file (to prevent loss of data if
// bench-routes goes down, acting like a copy of table buffer which can be re-read after a crash and restore the earlier
// state of table buffer) with rows in random order.
// The rows will be ordered only while flushing to the IO-buffer which is the present method as well.
func (tbuf *TableBuffer) flushToIOBuffer(flushIOBuffer bool) error {
	if tbuf.idata == 0 {
		// Return if there are no entries in the buffer.
		return nil
	}
	sort.SliceStable(tbuf.data, func(i, j int) bool {
		if tbuf.data[i].timestamp == 0 {
			return false
		}
		return tbuf.data[i].timestamp < tbuf.data[j].timestamp
	})
	for i := uint64(0); i < tbuf.idata; i++ {
		r := tbuf.data[i]
		if err := tbuf.writer.writeToTable(r.timestamp, r.id, r.valueSet); err != nil {
			return fmt.Errorf("flush to io-buffer: %w", err)
		}
	}
	tbuf.idata = 0
	tbuf.data = tbuf.data[:]
	if flushIOBuffer {
		if err := tbuf.writer.commit(); err != nil {
			return fmt.Errorf("flush to io-buffer: flushIOBuffer: %w", err)
		}
	}
	return nil
}
