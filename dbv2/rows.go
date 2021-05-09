package dbv2

import (
	"bufio"
	"fmt"
	"io"
)

const nullCell = "null"

var rowEmpty = row{-1, nullCell, nullCell}

type (
	row          struct {
		t int64
		rtype string
		v string
	}
	rowsIterator struct {
		rowLength   int // In number of bytes involved in a complete row.
		curr        uint64
		buf         []byte
		previousRow []byte
		reader      *bufio.Reader
	}
)

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

// readNext reads the next row. If returnRow is false, then the iterator just iterates to the next row, while returning
// an empty row. This is helpful when you just need to move down a line and avoid wasting computation resources behind
// parsing.
// If you are looking for last row, then the iterator uses a recovery mechanism to give back the previous row when
// it finds an EOF. This is usually the last row unless there is something fishy going on.
func (ri *rowsIterator) readNext(returnRow bool) (r row, isLast bool, err error) {
	// todo: unit test imp.
	ri.buf = ri.buf[:]
	_, err = ri.reader.Read(ri.buf)
	if err != nil {
		if err == io.EOF {
			isLast = true
			row, err := parseBytesToRow(ri.previousRow)
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
	if !returnRow {
		// This is helpful if the caller just wants to iterate over the rows and not return any stuff. This can be
		// done to move down a couple of lines and in such cases, we do not aim to do computation behind parsing
		// a complete row, which can save a lot of time and resources.
		return rowEmpty, false, nil
	}
	row, err := parseBytesToRow(ri.buf)
	if err != nil {
		return rowEmpty, false, fmt.Errorf("read-next: %w", err)
	}
	return row, false, nil
}
