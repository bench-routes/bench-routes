package dbv2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	filePath  = "testdata/test_writer_file"
	indexPath = "testdata/test_writer_file"
)

func TestTableReader(t *testing.T) {
	r, err := NewTableReader(filePath)
	require.NoError(t, err)

	ir := NewTableIndex(indexPath)
	require.NoError(t, ir.Load())

	reader := r.Reader()
	rows, err := reader.ReadAll(ir.Get(1))
	require.NoError(t, err)

	expectedRows := []row{
		{1, "1", "1", "234.5:11.2:3"},
		{10, "1", "1", "1.5:11.2:3"},
		{100, "1", "1", "3.5:11.2:3"},
		{200, "1", "1", "5.5:11.2:3"},
		{300, "1", "1", "7:11.2:3"},
		{500, "1", "1", "9.5:11.2:3"},
		{800, "1", "1", "11.5:11.2:3"},
	}

	require.Equal(t, expectedRows, *rows)

	// todo: return readers
	readerAlt := r.Reader()
	rows, err = readerAlt.ReadAll(ir.Get(2))
	require.NoError(t, err)

	expectedRows = []row{
		{100, "2", "1", "234.5:11.2:3"},
		{200, "2", "1", "234.5:11.2:3"},
		{300, "2", "1", "234.5:11.2:3"},
	}

	require.Equal(t, expectedRows, *rows)
}
