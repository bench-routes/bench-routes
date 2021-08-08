package querier

import (
	"testing"

	tsdb "github.com/bench-routes/bench-routes/tsdb"
)

var tests []string = []string{
	"./testdata/test_ping.json",
	"./testdata/test_jitter.json",
	"./testdata/test_monitor.json",
}

func TestQurier(t *testing.T) {
	for _, path := range tests {
		bstream, err := tsdb.FetchChainStream(path)
		if err != nil {
			t.Fatalf("loading chain error: %s\n", err)
		}
		sortedStream := sortStream(bstream)
		for i := 0; i < (len(*sortedStream) - 1); i++ {
			if (*sortedStream)[i].NormalizedTime > (*sortedStream)[i+1].NormalizedTime {
				t.Fatalf("error sorting : %d > %d \n", (*sortedStream)[i].NormalizedTime, (*sortedStream)[i+1].NormalizedTime)
			}
		}
		tIndex := len(*sortedStream)/2 + len(*sortedStream)/6

		gIndex := binSearch(sortedStream, (*sortedStream)[tIndex].NormalizedTime, 0, len(*sortedStream))

		if tIndex != gIndex {
			t.Fatalf("binSearch not working : %d != %d\n", tIndex, gIndex)
		}
		strIdx := len(*sortedStream) / 3
		endIdx := len(*sortedStream) / 2

		q, err := New(TypeRange, path, (*sortedStream)[strIdx].NormalizedTime, (*sortedStream)[endIdx].NormalizedTime)
		if err != nil {
			t.Fatalf("query error : %s\n", err.Error())
		}

		stream, err := q.fetchBlocks(sortedStream)

		if err != nil {
			t.Fatalf("fetching block error : %s\n", err.Error())
		}

		if len(stream) != (endIdx - strIdx + 1) {
			t.Fatalf("fetching block error : length of response is not expected : %d != %d\n", len(stream), endIdx-strIdx+1)
		}

		_, err = q.Exec()
		if err != nil {
			t.Fatalf("query exec error : %s\n", err.Error())
		}
	}
}
