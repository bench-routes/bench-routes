package tsdb

import (
	"fmt"
	"testing"
	"time"
)

var (
	chain = Chain{
		Path:           "../test-files/loadFromStorage_testdata/test1.json",
		LengthElements: 0,
		Chain:          []Block{},
		Size:           0,
	}
)

func TestInit(t *testing.T) {
	blocks := chain.Init()
	if len(blocks.Chain) == 0 {
		t.Errorf("tsdb Init not working as expected")
	} else {
		t.Log("printing block values ...")
		t.Log(blocks)
	}
}

func TestAppend(t *testing.T) {
	chain := chain.Init()
	b := Block{
		PrevBlock:      nil,
		NextBlock:      nil,
		Timestamp:      time.Now(),
		NormalizedTime: 1568705420,
		Datapoint:      20,
	}

	c := chain.Append(b)
	fmt.Println(c.Chain[len(c.Chain)-1].Datapoint)
	if c.Chain[len(c.Chain)-1].Datapoint == 20 {
		t.Logf("Block Append Successful")
	} else {
		t.Errorf("Block Append Unsuccessful")
	}
}

func TestPopPreviousNBlocks(t *testing.T) {
	chain := chain.Init()
	_, err := chain.PopPreviousNBlocks(10)
	if err != nil {
		t.Logf(err.Error())
	} else {
		t.Logf("Block removal worked properly")
	}
}

func TestGetPositionalPointerNormalized(t *testing.T) {
	chain := chain.Init()
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
	chain := chain.Init()
	if chain.Save() {
		t.Logf("tsdb Save works as expected")
	}
}

func TestChainTraversal(t *testing.T) {
	chain := chain.Init()
	c := chain.Chain
	node := c[0]
	if node.PrevBlock != nil {
		t.Errorf("corrupted chain")
	}
	count := 1
	for {
		node = *node.NextBlock
		if node.NextBlock == nil {
			break
		}
		count++
	}

	if count == len(c) {
		t.Errorf("corrupted chain")
	} else {
		t.Logf("Succesful traversal")
	}
}
