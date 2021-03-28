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
	valueSeparator = ':'
	typeSeparator  = '|'
	space          = " "
	newLineSymbol  = '\n'
	rowEndSymbol   = '\ufffe'
)

// TODO (harkishen): replace all fmt.Println()s to logfmt format.

var (
	validLengthString = strings.Repeat(space, 100)
	validMaxBytes     = []byte(validLengthString)
	// validLength corresponds to the maximum length of a line in a data-table.
	validLength        = len(validLengthString)
	numBytesSingleLine = len(validMaxBytes)
	bufferSize         = 10000
	writerBufferSize   = validLength * bufferSize
)

var (
	newLineSymbolByte = newLineSymbol
	rowEndSymbolByte  = rowEndSymbol
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
		if _, err = file.WriteString(fmt.Sprintf("maxValidLength: %d chars", validLength)); err != nil {
			os.Remove(file.Name())
			return dtbl, isCreated, fmt.Errorf("create data-table: writing maxValidLength: %w", err)
		}
		if _, err = file.Write([]byte("\n\n")); err != nil {
			os.Remove(file.Name())
			return dtbl, isCreated, fmt.Errorf("writing empty lines: %w", err)
		}
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
	dtbl.buffer = NewTableBuffer(dtbl, dtbl.mux, 200, 1000, time.Minute)
	return
}

type writeData struct {
	timestamp uint64
	id        uint64
	valueSet  string
}

// TableBuffer is a in-memory table that is in raw format and is yet to be committed.
type TableBuffer struct {
	bufferMux     sync.RWMutex
	data          []writeData
	cap           uint64
	size          uint64
	flushDeadline time.Duration
	tableCopy     *DataTable
	writer        *tableWriter
}

// NewTableBuffer returns a new TableBuffer.
func NewTableBuffer(dataTable *DataTable, mux *sync.RWMutex, cap uint64, ioBufferSize int, flushDeadline time.Duration) *TableBuffer {
	index := NewTableIndex(dataTable.path)
	seriesRelationIndex := NewSeriesRelation(dataTable.path)
	return &TableBuffer{
		cap:           cap,
		flushDeadline: flushDeadline,
		tableCopy:     dataTable,
		data:          make([]writeData, 0),
		writer:        newTableWriter(dataTable, mux, ioBufferSize, index, seriesRelationIndex),
	}
}

func (tbuf *TableBuffer) Write(timestamp, id uint64, valueSet string) error {
	tbuf.bufferMux.Lock()
	defer tbuf.bufferMux.Unlock()
	if tbuf.size > tbuf.cap {
		if err := tbuf.flushToIOBuffer(false); err != nil {
			return fmt.Errorf("table-buffer write: %w", err)
		}
	}
	tbuf.data = append(tbuf.data, writeData{timestamp: timestamp, id: id, valueSet: valueSet})
	tbuf.size++ // This does not need atomic operation as it is protected by the bufferMux lock by default.
	return nil
}

// flushToIOBuffer is not concurrent safe. The caller is expected to use a lock.
// TODO: this can be problematic in situations when the API takes more than flushDeadline. This can be resolved by
// increasing the flushDuration to an hour and maintain the state of table buffer in a file (to prevent loss of data if
// bench-routes goes down, acting like a copy of table buffer which can be re-read after a crash and restore the earlier
// state of table buffer) with rows in random order.
// The rows will be ordered only while flushing to the IO-buffer which is the present method as well.
func (tbuf *TableBuffer) flushToIOBuffer(flushIOBuffer bool) error {
	if tbuf.size == 0 {
		// Return if there are no entries in the buffer.
		return nil
	}
	sort.SliceStable(tbuf.data, func(i, j int) bool {
		return tbuf.data[i].timestamp < tbuf.data[j].timestamp
	})
	for i := uint64(0); i < tbuf.size; i++ {
		r := tbuf.data[i]
		if err := tbuf.writer.writeToTable(r.timestamp, r.id, r.valueSet); err != nil {
			return fmt.Errorf("flush to io-buffer: %w", err)
		}
	}
	tbuf.size = 0
	tbuf.data = tbuf.data[:]
	if flushIOBuffer {
		if err := tbuf.writer.commit(); err != nil {
			return fmt.Errorf("flush to io-buffer: flushIOBuffer: %w", err)
		}
	}
	return nil
}

type tableWriter struct {
	index               *tableIndex
	seriesRelIndex      *seriesRelation
	currentBytePosition uint64
	tableWriter         *bufio.Writer
	mux                 *sync.RWMutex
	minAcceptableTs     uint64
}

func newTableWriter(dataTable *DataTable, mux *sync.RWMutex, tableIOBufferSize int, index *tableIndex, seriesRelationIndex *seriesRelation) *tableWriter {
	return &tableWriter{
		mux:            mux,
		index:          index,
		seriesRelIndex: seriesRelationIndex,
		tableWriter:    bufio.NewWriterSize(dataTable, tableIOBufferSize*validLength), // Buffer size corresponds to total number of rows.
	}
}

func ConvertValueToValueSet(data ...string) string {
	//todo: needs a test and a probable benchmark in future.
	var s string
	for i := range data {
		s = fmt.Sprintf("%s%s%c", s, data[i], valueSeparator)
	}
	return s[:len(s)-1] // Ignore the last pipe.
}

// writeToTable writes to the table. It must get timestamps in increasing order only.
func (w *tableWriter) writeToTable(timestamp, id uint64, valueSet string) error {
	if timestamp < w.minAcceptableTs {
		return fmt.Errorf("timestamp cannot be less than minAcceptableTs: wanted >= %d received %d", w.minAcceptableTs, timestamp)
	}
	serialized, err := serializeWrite(timestamp, id, valueSet)
	if err != nil {
		return fmt.Errorf("serialize-write: %w", err)
	}
	w.mux.Lock()
	defer w.mux.Unlock()
	endAt, err := w.tableWriter.Write(serialized)
	if err != nil {
		if err == bufio.ErrBufferFull {
			if err = w.commit(); err != nil {
				return fmt.Errorf("commit: %w", err)
			}
			if _, err = w.tableWriter.Write(serialized); err != nil {
				return fmt.Errorf("failed retrying write single after flush: %w", err)
			}
		} else {
			return fmt.Errorf("tableWriter.Write(): %w", err)
		}
	}
	currPos := atomic.LoadUint64(&w.currentBytePosition)
	if err := w.index.Add(id, currPos); err != nil {
		fmt.Println("warn: ", "index.Add(): ", err.Error())
	}
	atomic.AddUint64(&w.currentBytePosition, uint64(endAt))
	if timestamp > w.minAcceptableTs {
		w.minAcceptableTs = timestamp
	}
	return nil
}

func (w *tableWriter) commit() error {
	w.mux.Lock()
	defer w.mux.Unlock()
	if err := w.index.Flush(); err != nil {
		return fmt.Errorf("commit: flushing index: %w", err)
	}
	if err := w.tableWriter.Flush(); err != nil {
		return fmt.Errorf("error occurred while flushing: %w", err)
	}
	return nil
}

func serializeWrite(timestamp, seriesID uint64, value string) ([]byte, error) {
	// todo: needs a test
	// todo: remove type info as we now have start indices. This needs to be updated everywhere
	str := fmt.Sprintf("%d%c%d%c%s\n", timestamp, typeSeparator, seriesID, typeSeparator, value)
	if l := len(str); l > validLength {
		return nil, fmt.Errorf("length greater than the valid-length: received %d wanted %d", l, validLength)
	}
	return []byte(str), nil
}
