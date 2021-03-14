package dbv2

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	space          = " "
	writeSeparator = "|"
	colon          = ":"
)

// TODO (harkishen): replace all fmt.Println()s to logfmt format.

var (
	validLengthString = strings.Repeat(space, 100)
	// validLength corresponds to the maximum length of a line in a data-table.
	validLength      = len(validLengthString)
	bufferSize       = 10000
	writerBufferSize = validLength * bufferSize
)

type DataTable struct {
	*os.File
	mux               *sync.RWMutex
	offset            int16
	minWriteTimestamp uint64
	tableBuffer       []*writeData
	writer            *TableWriter
}

func CreateDataTable(path string) (*DataTable, error) {
	var file *os.File
	// Verify if file exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("creating new file")
		file, err = os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("create data-table: creating new data-table file: %w", err)
		}
		if _, err = file.WriteString(fmt.Sprintf("maxValidLength: %d chars", validLength)); err != nil {
			os.Remove(file.Name())
			return nil, fmt.Errorf("create data-table: writing maxValidLength: %w", err)
		}
		if _, err = file.Write([]byte("\n\n")); err != nil {
			os.Remove(file.Name())
			return nil, fmt.Errorf("writing empty lines: %w", err)
		}
	} else {
		fmt.Println("opening new file")
		file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return nil, fmt.Errorf("create data-table: open existing data-table file: %w", err)
		}
		// TODO (harkishen): verify the maxLength by reading from file.
	}
	dtbl := &DataTable{
		File:              file,
		mux:               new(sync.RWMutex),
		offset:            3,
		minWriteTimestamp: 0,
	}
	return dtbl, nil
}

type tableWriter struct {
	tableWriter     *bufio.Writer
	mux             *sync.RWMutex
	minAcceptableTs uint64
}

type writeData struct {
	timestamp uint64
	id        uint64
	valueSet  string
}

var writeDataPool = sync.Pool{New: func() interface{} { return new(writeData) }}

type TableBuffer struct {
	bufferMux     sync.RWMutex
	data          []*writeData
	cap           uint64
	size          uint64
	flushDeadline time.Duration
	tableWriter   *tableWriter
}

func NewTableBuffer(dataTable *DataTable, mux *sync.RWMutex, cap uint64, ioBufferSize int, flushDeadline time.Duration) *TableBuffer {
	return &TableBuffer{
		cap:           cap,
		flushDeadline: flushDeadline,
		tableWriter:   newTableWriter(dataTable, mux, ioBufferSize),
	}
}

func (tbuf *TableBuffer) Write(timestamp, id uint64, valueSet string) {
	tbuf.bufferMux.Lock()
	if tbuf.size > tbuf.cap {
		//	 flush here
	}
	row := writeDataPool.Get().(*writeData)
	row.timestamp = timestamp
	row.id = id
	row.valueSet = valueSet
	tbuf.data = append(tbuf.data, row)
	atomic.AddUint64(&tbuf.size, 1)
}

// flushToIOBuffer is not concurrent safe. The caller is expected to use a lock.
func (tbuf *TableBuffer) flushToIOBuffer() error {
	if tbuf.size == 0 {
		// Return if there are no entries in the buffer.
		return nil
	}
	release := func(start uint64) {
		for i := start; i < tbuf.size; i++ {
			writeDataPool.Put(tbuf.data[i])
		}
	}
	sort.SliceStable(tbuf.data, func(i, j int) bool {
		return tbuf.data[i].timestamp < tbuf.data[j].timestamp
	})
	for i := uint64(0); i < tbuf.size; i++ {
		r := tbuf.data[i]
		if err := tbuf.tableWriter.writeToTable(r.timestamp, r.id, r.valueSet); err != nil {
			fmt.Println("write to table failed. putting rows back into pool")
			release(i)
			return fmt.Errorf("flush to io-buffer: %w", err)
		}
		writeDataPool.Put(tbuf.data[i])
	}
	return nil
}

func newTableWriter(dataTable *DataTable, mux *sync.RWMutex, tableIOBufferSize int) *tableWriter {
	return &tableWriter{
		tableWriter: bufio.NewWriterSize(dataTable, tableIOBufferSize*validLength), // Buffer size corresponds to total number of rows.
		mux:         mux,
	}
}

func ConvertValueToValueSet(data ...string) string {
	//todo: needs a test
	s := ""
	for i := range data {
		s += data[i] + colon
	}
	return s[:len(s)-1] // Ignore the last pipe.
}

func (w *tableWriter) writeToTable(timestamp, id uint64, valueSet string) error {
	serialized, err := serializeWrite(timestamp, id, valueSet)
	if err != nil {
		return fmt.Errorf("serialize-write: %w", err)
	}
	w.mux.Lock()
	defer w.mux.Unlock()
	if _, err = w.tableWriter.Write(serialized); err != nil {
		if err == bufio.ErrBufferFull {
			fmt.Println("buffer full, flushing to disk and attempting again")
			if err = w.commit(); err != nil {
				return fmt.Errorf("commit: %w", err)
			}
			if _, err = w.tableWriter.Write(serialized); err != nil {
				return fmt.Errorf("failed retrying write single after flush: %w", err)
			}
		}
	}
	return nil
}

func (w *tableWriter) commit() error {
	w.mux.Lock()
	defer w.mux.Unlock()
	if err := w.tableWriter.Flush(); err != nil {
		return fmt.Errorf("error occurred while flushing: %w", err)
	}
	return nil
}

func serializeWrite(timestamp, seriesID uint64, value string) ([]byte, error) {
	// todo: needs a test
	str := fmt.Sprintf("%d%s%d%s%s\n", timestamp, writeSeparator, seriesID, writeSeparator, value)
	if l := len(str); l > validLength {
		return nil, fmt.Errorf("length greater than the valid-length: received %d wanted %d", l, validLength)
	}
	return []byte(str), nil
}
