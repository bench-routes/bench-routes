package tsdb

import (
	"testing"
	"time"
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

func TestAppend(t *testing.T) {
	_, chain := chain.Init()
	b := Block{
		PrevBlock:      nil,
		NextBlock:      nil,
		Timestamp:      time.Now(),
		NormalizedTime: 1568705420,
		Datapoint:      20,
	}

	status, c := chain.Append(&b)
	if status {
		if c.lengthElements == chain.lengthElements+1 {
			t.Logf("Block Append Successful")
		} else {
			t.Errorf("Block Append Unsuccessful")
		}
	} else {
		t.Errorf("Block Append Unsuccessful")
	}
}

func TestPopPreviousNBlocks(t *testing.T) {
	_, chain := chain.Init()
	chain, err := chain.PopPreviousNBlocks(10)
	if err != nil {
		t.Logf(err.Error())
	} else {
		t.Logf("Block removal worked properly")
	}
}

func TestGetPositionalPointerNormalized(t *testing.T) {
	_, chain := chain.Init()
	var normalizedTime int64 = 1568705425
	block := *chain.GetPositionalPointerNormalized(1568705425)
	if block.NormalizedTime == normalizedTime {
		t.Log("Test Successful")
		t.Log(block)
	} else {
		t.Errorf("Block not found")
	}
}

func TestSave(t *testing.T) {
	_, chain := chain.Init()
	if chain.Save() {
		t.Logf("tsdb Save works as expected")
	}
}
