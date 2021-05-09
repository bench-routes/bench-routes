package dbv2

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
)

var (
	readerBufferSize = defaultBufferSize
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
		reader: bufio.NewReader(f),
	}
	tr.rowsItr = newRowsIterator(tr.reader, 0)
	return tr, nil
}

// Parse traverses the rows of the table.
func (dr *TableReader) Parse() error {
	var (
		r      row
		err    error
		isLast bool
	)
	for {
		r, isLast, err = dr.rowsItr.readNext(false)
		if err != nil {
			return fmt.Errorf("parse: %w", err)
		}
		if isLast {
			break
		}
	}
	dr.dataTable.minWriteTimestamp = uint64(r.t)
	return nil
}

type reader struct {
	tableReader *bufio.Reader
	rowsItr   *rowsIterator
}

func (dr *TableReader) Reader() *reader {
	r := bufio.NewReader(dr.dataTable.File)
	return &reader{r, newRowsIterator(r, 0)}
}

// skipBytes skips the given number bytes from the current position.
func (dr *reader) skipBytes(numBytes int) error {
	_, err := dr.tableReader.Discard(numBytes)
	if err != nil {
		return fmt.Errorf("skipBytes: %w", err)
	}
	return nil
}

func (dr *reader) ReadAll(seriesIdIndex []uint64) (*[]row, error) {
	var (
		skip, previousRead uint64
		bytesRead          int
		r                  row
		err                error
		result             []row
	)
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

// row returns the row after parsing the byte slice from current reader position.
func (dr *reader) row() (row, int, error) {
	line, err := dr.tableReader.ReadBytes(newLineSymbol)
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

func parseBytesToRow(buf []byte) (row, error) {
	var (
		row = buf

		timestampBuf []byte
		seriesIdBuf  []byte
		typeIdBuf    []byte
		valueBuf     []byte
	)
	// Parse timestamp.
	index := bytes.Index(row, []byte{typeSeparator})
	timestampBuf = row[:index]
	row = row[index+1:]
	fmt.Println("row now", string(row))

	// Parse series id.
	index = bytes.Index(row, []byte{typeSeparator})
	seriesIdBuf = row[:index]
	row = row[index+1:]

	// Parse type id.
	index = bytes.Index(row, []byte{typeSeparator})
	typeIdBuf = row[:index]
	row = row[index+1:]

	// Parse value.
	index = bytes.Index(row, []byte{typeSeparator})
	valueBuf = row

	timestampInt, err := strconv.Atoi(string(timestampBuf))
	if err != nil {
		return rowEmpty, fmt.Errorf("parse-row: corrupted timestamp byte slice: %s", string(timestampBuf))
	}

	return makeRow(int64(timestampInt), seriesIdBuf, typeIdBuf, valueBuf), nil
}

// makeRow builds a table row.
func makeRow(timestamp int64, seriesId, typeId, value []byte) row {
	return row{timestamp, string(seriesId), string(typeId), string(value)}
}
