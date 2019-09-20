package tsdb

import (
	"testing"
)

var (
	chain = Chain{
		path:           "../test-files/loadFromStorage_testdata/test1.json",
		lengthElements: 0,
		chain:          []Block{},
		size:           0,
	}
)

func TestInit(t *testing.T) {
	blocks, _ := chain.Init()
	if len(*blocks) == 0 {
		t.Errorf("tsdb Init not working as expected")
	} else {
		t.Log("printing block values ...")
		t.Log(blocks)
	}
}

func TestSave(t *testing.T) {
	_, chain := chain.Init()
	if chain.Save() {
		t.Logf("tsdb Save works as expected")
	}
}
