package dbv2

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateFile_Write_Read(t *testing.T) {
	contents := []struct {
		timestamp, seriesID uint64
		values              []string
	}{
		{
			timestamp: 100,
			seriesID:  2,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 200,
			seriesID:  2,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 300,
			seriesID:  2,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 1,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 10,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 100,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 200,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 300,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 500,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 800,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
	}
	expected := `1|1|234.5:11.2:3
10|1|234.5:11.2:3
100|2|234.5:11.2:3
100|1|234.5:11.2:3
200|2|234.5:11.2:3
200|1|234.5:11.2:3
300|2|234.5:11.2:3
300|1|234.5:11.2:3
500|1|234.5:11.2:3
800|1|234.5:11.2:3
`
	testFile := "test_writer_file"
	dtbl, created, err := OpenRWDataTable(testFile)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(testFile))
		require.NoError(t, os.Remove(testFile+".index"))
	}()
	require.True(t, created, "RWDatatable")

	for _, c := range contents {
		err := dtbl.buffer.Write(c.timestamp, c.seriesID, ConvertValueToValueSet(c.values[0], c.values[1], c.values[2]))
		require.NoError(t, err, "writing contents")
		// Keeping flushToIOBuffer(true) in for loop, is the actual tough part for the test to pass. Passing here would mean there is no
		// duplication of data (which is what we want). Ideally, this would be kept after the loop for write is done, but since this is testing,
		// we want to make sure for edge cases.
	}
	err = dtbl.buffer.flushToIOBuffer(true)
	require.NoError(t, err)
	bSlice, err := ioutil.ReadFile(testFile)
	require.NoError(t, err)
	require.Equal(t, expected, string(bSlice), "matching write result")
	fmt.Println(dtbl.buffer.writer.index)

	// Reading.
	reader, err := NewTableReader(testFile)
	require.NoError(t, err)
	rows, err := reader.read(dtbl.buffer.writer.index.Get(1))
	require.NoError(t, err)
	fmt.Println(rows)
}

func TestUnorderedInserts(t *testing.T) {
	contents := []struct {
		timestamp, seriesID uint64
		values              []string
	}{
		{
			timestamp: 1000,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 10,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 500,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 200,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 300000,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 700,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 1000,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 5001,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 20,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 3000,
			seriesID:  1,
			values:    []string{"234.5", "11.2", "3"},
		},
	}
	expected := `10|1|234.5:11.2:3
20|1|234.5:11.2:3
200|1|234.5:11.2:3
500|1|234.5:11.2:3
700|1|234.5:11.2:3
1000|1|234.5:11.2:3
1000|1|234.5:11.2:3
3000|1|234.5:11.2:3
5001|1|234.5:11.2:3
300000|1|234.5:11.2:3
`
	testFile := "test_unordered_inserts_file"
	dtbl, _, err := OpenRWDataTable(testFile)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(testFile))
	}()
	for _, c := range contents {
		err := dtbl.buffer.Write(c.timestamp, c.seriesID, ConvertValueToValueSet(c.values[0], c.values[1], c.values[2]))
		require.NoError(t, err, "writing contents")
	}
	err = dtbl.buffer.flushToIOBuffer(true)
	require.NoError(t, err)
	bSlice, err := ioutil.ReadFile(testFile)
	require.NoError(t, err)
	require.Equal(t, []byte(expected), bSlice, "matching write result")
	reader, err := NewTableReader(testFile)
	require.NoError(t, err)
	rows, err := reader.read(dtbl.buffer.writer.index.Get(1))
	require.NoError(t, err)
	fmt.Println(rows)
}
