package tsdb

import (
	"testing"
)

var (
	chain = Chain{
		Path:           "../tests/loadFromStorage_testdata/test1.json",
		LengthElements: 0,
		Chain:          []Block{},
		Size:           0,
	}
)

func TestInit(t *testing.T) {
	blocks := chain.Init()
	if len(blocks.Chain) == 0 {
		t.Errorf("tsdb Init not working as expected")
	}
}

func TestAppend(t *testing.T) {
	chain := chain.Init()
	b := *GetNewBlock("", "20")

	c := chain.Append(b)
	if c.Chain[len(c.Chain)-1].Datapoint != "20" {
		t.Errorf("Block Append Unsuccessful")
	}
}

func TestPopPreviousNBlocks(t *testing.T) {
	chain := chain.Init()
	_, err := chain.PopPreviousNBlocks(10)
	if err != nil {
		t.Logf(err.Error())
	}
}

func TestGetPositionalPointerNormalized(t *testing.T) {
	chain := chain.Init()
	var outOfRangeTime int64 = 21568705425
	var notFoundTime int64 = 1568705420

	_, errOutOfRange := chain.GetPositionalIndexNormalized(outOfRangeTime)
	if errOutOfRange == nil {
		t.Errorf("Out of range value check does not work")
	}

	_, errNotFound := chain.GetPositionalIndexNormalized(notFoundTime)
	if errNotFound == nil {
		t.Errorf("Check for element not found does not work")
	}
}

func TestSave(t *testing.T) {
	chain := chain.Init()
	chain.commit()
}

func TestChainSequence(t *testing.T) {
	chain := chain.Init()
	c := chain.Chain
	var last int64
	for i := 0; i < len(c)-1; i++ {
		if c[i].NormalizedTime > last {
			last = c[i].NormalizedTime
		} else {
			t.Errorf("Wrong chain sequence")
		}
	}
}
