package tsdb

import (
	"testing"
)

var (
	chain = Chain{
		Path:           "../tests/loadFromStorage_testdata/test1.json",
		LengthElements: 0,
		Chain:          []Block{},
	}
)

func TestInit(t *testing.T) {
	chain.Init()
}

func TestAppend(t *testing.T) {
	chain := chain.Init()
	b := *GetNewBlock("", "20")

	c := chain.Append(b)
	if c.Chain[len(c.Chain)-1].Datapoint != "20" {
		t.Errorf("Block Append Unsuccessful")
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
