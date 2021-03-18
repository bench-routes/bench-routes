package dbv2

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

var (
	readerBufferSize = validLength * (bufferSize / 10)
	// readerOffsetBytes represents the number of bytes that occupy the meta-data section of the table. This can be discarded
	// during parse operation as this is only to know the type of the table.
	readerOffsetBytes = len([]byte(`maxValidLength: 100 chars

`))
)

type Reader interface {
	Read()
}

type DataReader struct {
	dataTable *DataTable
	reader    *bufio.Reader
}

func NewDataReader(path string) (*DataReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("new data-reader: %w", err)
	}
	return &DataReader{
		dataTable: &DataTable{
			File: f,
		},
		reader: bufio.NewReaderSize(f, readerBufferSize),
	}, nil
}

// discardInitialOffset discards the initial meta-data of the table file.
func (dr *DataReader) discardInitialOffset() error {
	discarded, err := dr.reader.Discard(readerOffsetBytes)
	if err != nil {
		return fmt.Errorf("discard initial offset: %w", err)
	}
	if discarded != readerOffsetBytes {
		return fmt.Errorf("file does not have sufficient bytes written. Corrupted format")
	}
	return nil
}

func (dr *DataReader) Parse() error {
	if err := dr.discardInitialOffset(); err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	var (
		r      row
		err    error
		isLast bool
		itr    = newRowsIterator(dr.reader, numBytesSingleLine)
	)
	for {
		r, isLast, err = itr.readNext(true)
		if err != nil {
			return fmt.Errorf("parse: %w", err)
		}
		if isLast {
			fmt.Println("got to last row", r)
			break
		}
	}
	dr.dataTable.minWriteTimestamp = uint64(r.rowTimestamp)
	return nil
}

type (
	rowTimestamp int64
	rowType      string
	rowValue     string
	row          struct {
		rowTimestamp
		rowType
		rowValue
	}
)

const nullCell = "null"

var rowEmpty = row{-1, nullCell, nullCell}

func makeRow(timestamp rowTimestamp, rtype rowType, value rowValue) row {
	return row{timestamp, rtype, value}
}

type rowsIterator struct {
	rowLength   int // In number of bytes involved in a complete row.
	curr        uint64
	buf         []byte
	previousRow []byte
	reader      *bufio.Reader
}

// newRowsIterator returns a new rowsIterator. A rowsIterator must be used only on the row elements. This means that
// it must not be used from the initial metadata string, otherwise the operation will complain about corruptness. Hence,
// rowsIterator must be called only after the initial Parse() of *DataReader.
func newRowsIterator(reader *bufio.Reader, rowLength int) *rowsIterator {
	return &rowsIterator{
		buf:         make([]byte, rowLength),
		previousRow: make([]byte, rowLength),
		rowLength:   rowLength,
		reader:      reader,
	}
}

func parseRow(buf []byte) (row, error) {
	rowEndingPosition := bytes.Index(buf, []byte{rowEndSymbolByte})
	row := buf[:rowEndingPosition]

	timestampEndIndex := bytes.Index(row, []byte{typeSeparator})
	typeEndIndex := bytes.LastIndex(row, []byte{typeSeparator})
	timestampInt, err := strconv.Atoi(string(row[:timestampEndIndex]))
	if err != nil {
		return rowEmpty, fmt.Errorf("parse-row: corrupted timestamp byte slice: %s", string(row[:timestampEndIndex]))
	}

	timestamp := rowTimestamp(timestampInt)
	typeVal := rowType(row[timestampEndIndex+1 : typeEndIndex])
	value := rowValue(row[typeEndIndex+1:])

	return makeRow(timestamp, typeVal, value), nil
}

// readNext reads the next row. If returnRow is false, then the iterator just iterates to the next row, while returning
// an empty row. This is helpful when you just need to move down a line and avoid wasting computation resources behind
// parsing.
// If you are looking for last row, then the iterator uses a recovery mechanism to give back the previous row when
// it finds an EOF. This is usually the last row unless there is something fishy going on.
func (ri *rowsIterator) readNext(returnRow bool) (r row, isLast bool, err error) {
	// todo: unit test imp.
	ri.buf = ri.buf[:]
	n, err := ri.reader.Read(ri.buf)
	if err != nil {
		if err == io.EOF {
			isLast = true
			row, err := parseRow(ri.previousRow)
			if err != nil {
				return rowEmpty, false, fmt.Errorf("read-next: %w", err)
			}
			return row, true, nil
		}
		return rowEmpty, false, fmt.Errorf("read-next: %w", err)
	}
	copy(ri.previousRow, ri.buf)
	if len(string(ri.buf)) == 0 {
		return rowEmpty, false, nil
	}
	if n != numBytesSingleLine {
		return rowEmpty, false, fmt.Errorf("mismatch read: row can be corrupted: %s", string(ri.buf))
	}
	if !returnRow {
		// This is helpful if the caller just wants to iterate over the rows and not return any stuff. This can be
		// done to move down a couple of lines and in such cases, we do not aim to do computation behind parsing
		// a complete row, which can save a lot of time and resources.
		return rowEmpty, false, nil
	}
	row, err := parseRow(ri.buf)
	if err != nil {
		return rowEmpty, false, fmt.Errorf("read-next: %w", err)
	}
	return row, false, nil
}

// jumpNLines jumps n lines from the current position. The jump should never be greater than the maximum number of
// lines in the file.
func (ri *rowsIterator) jumpNLines(n uint64) error {
	// todo: test
	toBeDiscarded := ri.rowLength * int(n)
	discarded, err := ri.reader.Discard(toBeDiscarded)
	if err != nil {
		return fmt.Errorf("discard initial offset: %w", err)
	}
	if discarded != toBeDiscarded {
		return fmt.Errorf("jumNLines: lines are not properly aligned. Some rows can be corrupted")
	}
	return nil
}
