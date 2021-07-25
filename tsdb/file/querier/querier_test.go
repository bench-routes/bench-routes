package querier

import (
	"fmt"
	"testing"

	tsdb "github.com/bench-routes/bench-routes/tsdb/file"
)

func TestQurier(t *testing.T){
	bstream, err := tsdb.FetchChainStream("./testdata/test1.json")
	if err != nil {
		t.Fatalf("loading chain error: %s\n", err)
	}
	sortedStream := sortStream(bstream)
	for i:=0;i<(len(*sortedStream)-1);i++ {
		if (*sortedStream)[i].NormalizedTime > (*sortedStream)[i+1].NormalizedTime{
			t.Fatalf("error sorting : %d > %d \n",(*sortedStream)[i].NormalizedTime,(*sortedStream)[i+1].NormalizedTime)
		}
	}
	tIndex := len(*sortedStream)/2 + len(*sortedStream)/6

	gIndex := binSearch(sortedStream,(*sortedStream)[tIndex].NormalizedTime,0,len(*sortedStream))

	if tIndex != gIndex{
		t.Fatalf("binSearch not working : %d != %d\n",tIndex,gIndex);
	}
	strIdx := len(*sortedStream)/3
	endIdx := len(*sortedStream)/2

	q,err := New(TypeRange,"./testdata/test1.json",(*sortedStream)[strIdx].NormalizedTime,(*sortedStream)[endIdx].NormalizedTime)
	if err != nil {
		t.Fatalf("query error : %s\n",err.Error())
	}

	stream,err := q.fetchBlocks(sortedStream)

	if err != nil {
		t.Fatalf("fetching block error : %s\n",err.Error())
	}

	if len(stream) != (endIdx - strIdx +1) {
		t.Fatalf("fetching block error : length of response is not expected : %d != %d\n",len(stream),endIdx-strIdx+1)
	}

	res,err := q.Exec()
	if err != nil {
		t.Fatalf("query exec error : %s\n",err.Error())
	}
	fmt.Printf("%+v\n",res)
}