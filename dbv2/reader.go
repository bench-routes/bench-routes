package dbv2

import (
	"bufio"
	"bytes"
	"fmt"
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

// TableReader is a data-table file reader.
type TableReader struct {
	dataTable *DataTable
	reader    *bufio.Reader
	rowsItr   *rowsIterator
}

// NewTableReader returns a new TableReader.
func NewTableReader(path string) (*TableReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("new data-reader: %w", err)
	}
	tr := &TableReader{
		dataTable: &DataTable{
			File: f,
		},
		reader: bufio.NewReaderSize(f, readerBufferSize),
	}
	tr.rowsItr = newRowsIterator(tr.reader, numBytesSingleLine)
	return tr, nil
}

// discardInitialOffset discards the initial meta-data of the table file.
func (dr *TableReader) discardInitialOffset() error {
	discarded, err := dr.reader.Discard(readerOffsetBytes)
	if err != nil {
		return fmt.Errorf("discard initial offset: %w", err)
	}
	if discarded != readerOffsetBytes {
		return fmt.Errorf("file does not have sufficient bytes written. Corrupted format")
	}
	return nil
}

// Parse traverses the rows of the table.
func (dr *TableReader) Parse() error {
	if err := dr.discardInitialOffset(); err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	var (
		r      row
		err    error
		isLast bool
	)
	for {
		r, isLast, err = dr.rowsItr.readNext(true)
		if err != nil {
			return fmt.Errorf("parse: %w", err)
		}
		if isLast {
			break
		}
	}
	dr.dataTable.minWriteTimestamp = uint64(r.rowTimestamp)
	return nil
}

// skipBytes skips the given number bytes from the current position.
func (dr *TableReader) skipBytes(numBytes int) error {
	_, err := dr.reader.Discard(numBytes)
	if err != nil {
		return fmt.Errorf("skipBytes: %w", err)
	}
	return nil
}

// row returns the row after parsing the byte slice from current reader position.
func (dr *TableReader) row() (row, int, error) {
	line, err := dr.reader.ReadBytes(newLineSymbol)
	if err != nil {
		return rowEmpty, 0, fmt.Errorf("tableReader.row: %w", err)
	}
	bytesRead := len(line)
	row, err := parseBytesToRow(line[:len(line)-1]) // 1 corresponds to newLineSymbolByte. Ignore newLine symbol, otherwise the value field will contain a new line.
	if err != nil {
		return rowEmpty, bytesRead, fmt.Errorf("tableReader.row: %w", err)
	}
	return row, bytesRead, nil
}

func (dr *TableReader) read(seriesIdIndex []uint64) (*[]row, error) {
	dr.reader.Reset(dr.dataTable)
	var (
		skip, previousRead uint64
		bytesRead          int
		r                  row
		err                error
		result             []row
	)
	if err = dr.discardInitialOffset(); err != nil {
		return nil, fmt.Errorf("tableReader.read: %w", err)
	}
	for i := range seriesIdIndex {
		skip = seriesIdIndex[i]
		if i > 0 {
			skip -= seriesIdIndex[i-1]
		}
		skip -= previousRead
		if err = dr.skipBytes(int(skip)); err != nil {
			return nil, fmt.Errorf("tableReader.read: %w", err)
		}
		r, bytesRead, err = dr.row()
		if err != nil {
			return nil, fmt.Errorf("tableReader.read: %w", err)
		}
		previousRead = uint64(bytesRead)
		result = append(result, r)
	}
	return &result, nil
}

func parseBytesToRow(buf []byte) (row, error) {
	row := buf

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

// JumpNLines jumps n lines from the starting of the file, after the meta-data.
func (dr *TableReader) JumpNLines(n uint64) error {
	dr.reader.Reset(dr.dataTable)
	if err := dr.discardInitialOffset(); err != nil {
		return fmt.Errorf("tableRader.JumpNLines: %w", err)
	}
	if err := dr.rowsItr.jumpNLines(n); err != nil {
		return fmt.Errorf("tableReader.JumpNLines: %w", err)
	}
	return nil
}

// JumpNextNLines jumps n lines from the current position of the file reader.
func (dr *TableReader) JumpNextNLines(n uint64) error {
	if err := dr.rowsItr.jumpNLines(n); err != nil {
		return fmt.Errorf("tableReader.JumpNextNLines: %w", err)
	}
	return nil
}
