package dbv2

import (
	"bufio"
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	valueSeparator = ':'
	typeSeparator  = '|'
	space          = " "
	newLineSymbol  = '\n'
	colon = valueSeparator
	rowEndSymbol   = '\ufffe'
)

var (
	bufferSize         = 10000
	newLineSymbolByte = newLineSymbol
	rowEndSymbolByte  = rowEndSymbol
)

type tableWriter struct {
	index               *tableIndex
	currentBytePosition uint64
	tableWriter         *bufio.Writer
	mux                 *sync.RWMutex
	minAcceptableTs     uint64
}

func newTableWriter(dataTable *DataTable, mux *sync.RWMutex, tableIOBufferSize int) *tableWriter {
	return &tableWriter{
		mux:         mux,
		index: NewTableIndex(dataTable.path),
		tableWriter: bufio.NewWriterSize(dataTable, tableIOBufferSize), // Buffer size corresponds to total number of rows.
	}
}

// writeToTable writes to the table. It must get timestamps in increasing order only.
func (w *tableWriter) writeToTable(timestamp, id uint64, valueSet string) error {
	w.mux.Lock()
	defer w.mux.Unlock()
	if timestamp < w.minAcceptableTs {
		return fmt.Errorf("timestamp cannot be less than minAcceptableTs: wanted >= %d received %d", w.minAcceptableTs, timestamp)
	}
	serialized, err := serializeWrite(timestamp, id, valueSet)
	if err != nil {
		return fmt.Errorf("serialize-write: %w", err)
	}
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
	str := fmt.Sprintf("%d%c%d%c%s\n", timestamp, typeSeparator, seriesID, typeSeparator, value)
	return []byte(str), nil
}
