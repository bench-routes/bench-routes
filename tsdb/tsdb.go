package tsdb

import (
	"unsafe"
	"log"
	"time"
)

// Block use case block for the TSDB chain
type Block struct {
	prevBlock 		*Block
	nextBlock 		*Block
	datapoint 		int64
	normalizedTime 	uint64
	timestamp 		time.Time
}

// Chain contains Blocks arranged as a chain
type Chain struct {
	path 			string
	chain 			[]Block
	lengthElements 	uint64
	size 		   	uintptr
}

// TSDB implements the idea of tsdb
type TSDB interface {
	// Init helps to initialize the tsdb chain for the respective component. This function
	// should be capable to detect existing wals of the required type and build from the
	// local storage at the init of main thread and return the chain address in order to
	// have a minimal effect on the performance.
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

	// GetChain returns the address of chain.
	GetChain() *[]Block

	// Save saves or commits the chain in storage and returns success status.
	Save() bool
}

// Init initialize Chain properties
func (c Chain) Init(path *string) *[]Block {
	_, e := parse(*path)
	if e != nil {
		log.Printf("chain not found at %s. creating one ...", *path)
		c.path = *path
		c.lengthElements = 0
		c.size = unsafe.Sizeof(c)
		c.chain = []Block{}
		return &c.chain
	}
	c.path = *path
	return nil
}

