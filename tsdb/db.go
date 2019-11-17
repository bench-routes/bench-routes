package tsdb

import (
	"errors"
	"log"
	"sync"
	"time"
	"unsafe"
)

var (
	parser Parser
	// PingDBNames contains the name of the database corresponding to the unique config url
	PingDBNames = make(map[string]string)
	// FloodPingDBNames contains the name of the flood ping corresponding to the unique config url
	FloodPingDBNames = make(map[string]string)
	// GlobalPingChain contains chains of all the pings operating in bench-routes which has to be globally accessed
	// This is necessary as it helps to retain the parent values which are required for concurreny
	GlobalPingChain []*ChainPing

	// GlobalFloodPingChain contains chains of flood ping operations in bench-routes which has to be globally accessed
	// This is necessary as it helps to retain the parent values which are required for concurreny
	GlobalFloodPingChain []*ChainFloodPing
	//GlobalChain asdfafds
	GlobalChain []*Chain
	// GlobalResponseLength contains length of all the responses from a route
	GlobalResponseLength []*Chain
	// GlobalResponseDelay contains the all the delays in response from a request sent from a route
	GlobalResponseDelay []*Chain
	// GlobalResponseStatusCode contains the status code when a request is sent from to a route
	GlobalResponseStatusCode []*Chain
)

const (
	// BlockSeparation sets a separator for block datavalue
	BlockSeparation = "|"
)

// Block use case block for the TSDB chain
type Block struct {
	Datapoint      string // complex data would be decoded by using a blockSeparator
	NormalizedTime int64  // based on time.Unixnano()
	Type           string // would be used to decide the marshalling struct
	Timestamp      string
}

// BlockJSON helps refer Block as JSON
type BlockJSON struct {
	Datapoint      float32   `json:"datapoint"`
	NormalizedTime int64     `json:"normalizedTime"`
	Timestamp      time.Time `json:"timestamp"`
	Type           string    `json:"type"`
}

// Chain contains Blocks arranged as a chain
type Chain struct {
	Path           string
	Chain          []Block
	LengthElements int
	Size           uintptr
	mux            sync.Mutex
}

// TSDB implements the idea of tsdb
type TSDB interface {
	// Init helps to initialize the tsdb chain for the respective component. This function
	// should be capable to detect existing wals(write ahead log) of the required type and
	// build from the local storage at the init of main thread and return the chain address
	// in order to have a minimal effect on the performance.
	// Takes *path* as path to the existing chain or for creating a new one.
	// Returns address of the chain in RAM.
	Init() (*[]Block, Chain)

	// Append appends a new tsdb block passed as params to the most recent location (or
	// the last location) of the chain. Returns success status.
	Append(b Block) bool

	// GetPositionalPointerNormalized accepts the normalized time, searches for the block with that time
	// using jump search, and returns the address of the block having the specified normalized
	// time.
	GetPositionalPointerNormalized(n int64) (*Block, error)

	// PopPreviousNBlocks pops or removes **n** previous blocks from the chain and returns
	// success status.
	PopPreviousNBlocks(n uint64) (Chain, error)

	// GetChain returns the positional pointer address of the first element of the chain.
	GetChain() *[]Block

	// Save saves or commits the chain in storage and returns success status.
	Save() bool
}

// Init initialize Chain properties
func (c *Chain) Init() *Chain {
	c.mux.Lock()
	defer c.mux.Unlock()
	res, e := parse(c.Path)
	if e != nil {
		log.Printf("chain not found at %s. creating one ...", c.Path)
		c.LengthElements = 0
		c.Size = unsafe.Sizeof(c)
		c.Chain = []Block{}
		return c
	}

	raw := loadFromStorage(res)
	c.Chain = *formLinkedChainFromRawBlock(raw)
	c.LengthElements = len(c.Chain)
	c.Size = unsafe.Sizeof(c)
	return c
}

// Append function appends the new block in the chain
func (c *Chain) Append(b Block) *Chain {

	c.mux.Lock()
	defer c.mux.Unlock()
	c.Chain = append(c.Chain, b)
	c.Size = unsafe.Sizeof(c)
	c.LengthElements = len(c.Chain)
	l := c.LengthElements
	return c

}

// PopPreviousNBlocks pops last n elements from chain.
func (c *Chain) PopPreviousNBlocks(n int) (*Chain, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.LengthElements = len(c.Chain)
	l := c.LengthElements
	c.Chain = c.Chain[:len(c.Chain)-n]
	c.LengthElements = l - n
	c.Size = unsafe.Sizeof(c)
	return c, nil
}

// Save saves or commits the existing chain in the secondary memory.
// Returns the success status
func (c *Chain) Save() bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	log.Printf("writing chain of length %d", len(c.Chain))
	bytes := parser.ParseToJSON(c.Chain)
	e := saveToHDD(c.Path, bytes)
	if e != nil {
		panic(e)
	}
	return true
}

// GetPositionalPointerNormalized Returns block by searching the chain for the NormalizedTime
func (c *Chain) GetPositionalPointerNormalized(n int64) (*Block, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.LengthElements = len(c.Chain)
	if c.Chain[c.LengthElements-1].NormalizedTime < n || c.Chain[0].NormalizedTime > n {
		return nil, errors.New("Normalized time not in Chain range")
	}
	l, r, m := 0, c.LengthElements, 0

	for l <= r {
		m = l + (r-l)/2
		if c.Chain[m].NormalizedTime == n {
			return &c.Chain[m], nil
		}

		if c.Chain[m].NormalizedTime < n {
			l = m + 1
		} else {
			r = m - 1
		}
	}

	return nil, errors.New("Normalized time not found in chain")
}
