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
	var outOfRangeTime int64 = 21568705425
	var notFoundTime int64 = 1568705420

	_, errOutOfRange := chain.GetPositionalPointerNormalized(outOfRangeTime)
	if errOutOfRange != nil {
		t.Logf("Out of range value check works")
	}

	_, errNotFound := chain.GetPositionalPointerNormalized(notFoundTime)
	if errNotFound != nil {
		t.Logf("Check for element not found works")
	}

	block, _ := chain.GetPositionalPointerNormalized(normalizedTime)
	x := *block
	if x.NormalizedTime == normalizedTime {
		t.Logf("Block found")
		t.Log(x)
	}
}

func TestSave(t *testing.T) {
	_, chain := chain.Init()
	if chain.Save() {
		t.Logf("tsdb Save works as expected")
	}
}
