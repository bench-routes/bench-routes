package tsdb

import (
	"testing"
)

func TestInit(t *testing.T) {
	var (
		chain Chain
		path = "../test-files/loadFromStorage_testdata/test1.json"
	)
	blocks := *chain.Init(&path)
	if len(blocks) == 0 {
		t.Errorf("tsdb Init not working as expected")
	} else {
		t.Log("printing block values ...")
		t.Log(blocks)
	}
}