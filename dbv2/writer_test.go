package dbv2

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateFileAndWriteIt(t *testing.T) {
	contents := []struct {
		timestamp, seriesID uint64
		values              []string
	}{
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
	}
	expected := `maxValidLength: 100 chars

1|1|234.5:11.2:3
10|1|234.5:11.2:3
100|1|234.5:11.2:3
200|1|234.5:11.2:3
300|1|234.5:11.2:3
`
	testFile := "test_writer_file"
	dtbl, err := CreateDataTable(testFile)
	require.NoError(t, err)

	tblWriter := NewTableWriter(dtbl)
	for _, c := range contents {
		err := tblWriter.Write(c.timestamp, c.seriesID, ConvertValueToValueSet(c.values[0], c.values[1], c.values[2]))
		require.NoError(t, err, "writing contents")
		err = tblWriter.commit()
		require.NoError(t, err)
	}
	bSlice, err := ioutil.ReadFile(testFile)
	require.NoError(t, err)
	require.Equal(t, []byte(expected), bSlice, "matching write result")
	require.NoError(t, os.Remove(testFile))
}
