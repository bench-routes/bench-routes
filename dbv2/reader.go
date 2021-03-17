package dbv2

import (
	"bufio"
	"fmt"
	"os"
)

var (
	readerBufferSize = validLength * (bufferSize/10)
	readerOffsetBytes = len([]byte(`maxValidLength: 100 chars

`))
)

type Reader interface {
	Read()
}

type DataReader struct {
	dataTable *DataTable
	reader *bufio.Reader
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

func (dr *DataReader) Parse() {
	discarded, err := dr.reader.Discard(readerOffsetBytes)
	if err != nil {
		panic(err)
	}
	if discarded != readerOffsetBytes {
		panic("discarded must be same as offsetbytes")
	}
	var readLineBuf = make([]byte, numBytesSingleLine)
	n, err := dr.reader.Read(readLineBuf)
	if err != nil {
		panic(err)
	}
	if n < numBytesSingleLine {
		fmt.Println("short read")
	}
	fmt.Println("read was ", string(readLineBuf))
}


