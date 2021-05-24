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
		timestamp, seriesID, typeID uint64
		values                      []string
	}{
		{
			timestamp: 100,
			seriesID:  2,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 200,
			seriesID:  2,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 300,
			seriesID:  2,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 1,
			seriesID:  1,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 10,
			seriesID:  1,
			typeID:    1,
			values:    []string{"1.5", "11.2", "3"},
		},
		{
			timestamp: 100,
			seriesID:  1,
			typeID:    1,
			values:    []string{"3.5", "11.2", "3"},
		},
		{
			timestamp: 200,
			seriesID:  1,
			typeID:    1,
			values:    []string{"5.5", "11.2", "3"},
		},
		{
			timestamp: 300,
			seriesID:  1,
			typeID:    1,
			values:    []string{"7", "11.2", "3"},
		},
		{
			timestamp: 500,
			seriesID:  1,
			typeID:    1,
			values:    []string{"9.5", "11.2", "3"},
		},
		{
			timestamp: 800,
			seriesID:  1,
			typeID:    1,
			values:    []string{"11.5", "11.2", "3"},
		},
	}
	expected := `1|1|1|234.5:11.2:3
10|1|1|1.5:11.2:3
100|2|1|234.5:11.2:3
100|1|1|3.5:11.2:3
200|2|1|234.5:11.2:3
200|1|1|5.5:11.2:3
300|2|1|234.5:11.2:3
300|1|1|7:11.2:3
500|1|1|9.5:11.2:3
800|1|1|11.5:11.2:3
`
	testFile := "test_writer_file"
	dtbl, created, err := OpenRWDataTable(testFile)
	require.NoError(t, err)
	defer func() {
		//require.NoError(t, os.Remove(testFile))
		//require.NoError(t, os.Remove(testFile+".index"))
	}()
	require.True(t, created, "RWDatatable")

	for _, c := range contents {
		err := dtbl.buffer.Write(c.timestamp, c.seriesID, c.typeID, ConvertValueToValueSet(c.values[0], c.values[1], c.values[2]))
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
	r, err := NewTableReader(testFile)
	require.NoError(t, err)
	reader := r.Reader()
	rows, err := reader.ReadAll(dtbl.buffer.writer.index.Get(1))
	require.NoError(t, err)
	fmt.Println(rows)
}

func TestUnorderedInserts(t *testing.T) {
	contents := []struct {
		timestamp, seriesID, typeID uint64
		values                      []string
	}{
		{
			timestamp: 1000,
			seriesID:  1,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 10,
			seriesID:  1,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 500,
			seriesID:  1,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 200,
			seriesID:  1,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 300000,
			seriesID:  1,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 700,
			seriesID:  1,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 1000,
			seriesID:  1,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 5001,
			seriesID:  1,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 20,
			seriesID:  1,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
		{
			timestamp: 3000,
			seriesID:  1,
			typeID:    1,
			values:    []string{"234.5", "11.2", "3"},
		},
	}
	expected := `10|1|1|234.5:11.2:3
20|1|1|234.5:11.2:3
200|1|1|234.5:11.2:3
500|1|1|234.5:11.2:3
700|1|1|234.5:11.2:3
1000|1|1|234.5:11.2:3
1000|1|1|234.5:11.2:3
3000|1|1|234.5:11.2:3
5001|1|1|234.5:11.2:3
300000|1|1|234.5:11.2:3
`
	testFile := "test_unordered_inserts_file"
	dtbl, _, err := OpenRWDataTable(testFile)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(testFile))
		require.NoError(t, os.Remove(testFile+"."+indexFileExtension))
	}()
	for _, c := range contents {
		err := dtbl.buffer.Write(c.timestamp, c.seriesID, c.typeID, ConvertValueToValueSet(c.values[0], c.values[1], c.values[2]))
		require.NoError(t, err, "writing contents")
	}
	err = dtbl.buffer.flushToIOBuffer(true)
	require.NoError(t, err)
	bSlice, err := ioutil.ReadFile(testFile)
	require.NoError(t, err)
	require.Equal(t, []byte(expected), bSlice, "matching write result")
	r, err := NewTableReader(testFile)
	require.NoError(t, err)
	reader := r.Reader()
	rows, err := reader.ReadAll(dtbl.buffer.writer.index.Get(1))
	require.NoError(t, err)
	fmt.Println(rows)
}
