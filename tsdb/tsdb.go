package tsdb

import (
	"fmt"
	"log"
	"time"
	"unsafe"
)

// Block use case block for the TSDB chain
type Block struct {
	PrevBlock      *Block
	NextBlock      *Block
	Datapoint      int
	NormalizedTime int64
	Timestamp      time.Time
}

// BlockJSON helps reffer Block as JSON
type BlockJSON struct {
	Datapoint      int       `json:"datapoint"`
	NormalizedTime int64     `json:"normalizedTime"`
	Timestamp      time.Time `json:"timestamp"`
}

// Chain contains Blocks arranged as a chain
type Chain struct {
	path           string
	chain          []Block
	lengthElements int
	size           uintptr
}

// TSDB implements the idea of tsdb
type TSDB interface {
	// Init helps to initialize the tsdb chain for the respective component. This function
	// should be capable to detect existing wals(write ahead log) of the required type and
	// build from the local storage at the init of main thread and return the chain address
	// in order to have a minimal effect on the performance.
	// Takes *path* as path to the existing chain or for creating a new one.
	// Returns address of the chain in RAM.
	Init(path *string) *[]Block

	// Append appends a new tsdb block passed as params to the most recent location (or
	// the last location) of the chain. Returns success status.
	Append(b *Block) bool

	// GetPositionalPointer accepts the normalized time, searches for the block with that time
	// using jump search, and returns the address of the block having the specified normalized
	// time.
	GetPositionalPointer(n uint64) *Block

	// PopPreviousNBlocks pops or removes **n** previous blocks from the chain and returns
	// success status.
	PopPreviousNBlocks(n uint64) bool

	// GetChain returns the positional pointer address of the first element of the chain.
	GetChain() *[]Block

	// Save saves or commits the chain in storage and returns success status.
	Save() bool
}

// Init initialize Chain properties
func (c Chain) Init(path *string) *[]Block {
	res, e := parse(*path)
	// fmt.Println(res)
	if e != nil {
		log.Printf("chain not found at %s. creating one ...", *path)
		c.path = *path
		c.lengthElements = 0
		c.size = unsafe.Sizeof(c)
		c.chain = []Block{}
		return &c.chain
	}

	c.path = *path
	// fmt.Println(c.path)
	raw := loadFromStorage(res)
	// fmt.Println("rawwwwwwwwwwww::::::")
	// fmt.Println(raw)
	c.chain = *formLinkedChainFromRawBlock(raw)
	// fmt.Println("c.chain::::::;;;")
	// fmt.Println(c.chain)
	c.lengthElements = len(c.chain)
	c.size = unsafe.Sizeof(c)
	return &c.chain
}

// func main() {
// 	var (
// 		chain Chain
// 		path = "../test-files/loadFromStorage_testdata/test1.json"
// 	)
// 	blocks := *chain.Init(&path)
// 	if len(blocks) == 0 {
// 		fmt.Println("tsdb Init not working as expected")
// 	} else {
// 		fmt.Println("printing block values ...")
// 		fmt.Println(blocks)
// 	}
// }

func formLinkedChainFromRawBlock(a *[]BlockJSON) *[]Block {
	r := *a
	l := len(r)
	arr := []Block{}
	for i := 0; i < l; i++ {
		inst := Block{
			PrevBlock:      nil,
			NextBlock:      nil,
			Timestamp:      r[i].Timestamp,
			NormalizedTime: r[i].NormalizedTime,
			Datapoint:      r[i].Datapoint,
		}
		arr = append(arr, inst)
	}

	// form doubly linked list
	for i := 0; i < l; i++ {
		if i == 0 {
			arr[i].PrevBlock = nil
		} else {
			arr[i].PrevBlock = &arr[i-1]
		}
		if i == l-1 {
			arr[i].NextBlock = nil
		} else {
			arr[i].NextBlock = &arr[i+1]
		}
	}
	return &arr
}

func (c Chain) Append(b *Block) (bool, Chain) {

	c.chain = append(c.chain, *b)
	c.size = unsafe.Sizeof(c)
	c.lengthElements = len(c.chain)
	l := c.lengthElements
	c.chain[l-2].NextBlock = &c.chain[l-1]
	c.chain[l-1].PrevBlock = &c.chain[l-2]
	c.chain[l-1].NextBlock = nil
	return true, c

}
