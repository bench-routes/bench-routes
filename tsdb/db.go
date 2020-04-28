package tsdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/zairza-cetb/bench-routes/src/lib/logger"
)

const (
	// BlockDataSeparator sets a separator for block datavalue
	BlockDataSeparator = "|"
)

// Block use case block for the TSDB chain
type Block struct {
	Datapoint      string `json:"datapoint"`       // complex data would be decoded by using a blockSeparator
	NormalizedTime int64  `json:"normalized-time"` // based on time.Unixnano()
	Type           string `json:"type"`            // would be used to decide the marshalling struct
	Timestamp      string `json:"timestamp"`
}

// GetNewBlock creates and returns a new block with the specified type.
func GetNewBlock(blockType, value string) *Block {
	return &Block{
		Timestamp:      GetTimeStampCalc(),
		NormalizedTime: GetNormalizedTimeCalc(),
		Datapoint:      value,
		Type:           blockType,
	}
}

// Encode decodes the structure and marshals into a string
func (b Block) Encode() string {
	logger.File("decoding block type"+b.Type+" normalized as "+strconv.FormatInt(b.NormalizedTime, 10), "p")
	bbyte, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	return string(bbyte)
}

// GetType returns the type of the block
func (b Block) GetType() string {
	return b.Type
}

// GetDatapointEnc returns the data point to the caller.
// The encoded refers to the combined _(containing *|*)_ values in the string
// form.
func (b Block) GetDatapointEnc() string {
	return b.Datapoint
}

// GetNormalizedTimeStringified returns the normalized time of the block.
func (b Block) GetNormalizedTimeStringified() string {
	return string(b.NormalizedTime)
}

// GetNormalizedTime returns the normalized time of the block.
func (b Block) GetNormalizedTime() int64 {
	return b.NormalizedTime
}

// GetTimeStamp returns the timestamp of the block.
func (b Block) GetTimeStamp() string {
	return b.Timestamp
}

// Chain contains Blocks arranged as a chain
type Chain struct {
	Path           string
	Chain          []Block
	LengthElements int
	Size           uintptr
	mux            sync.Mutex
}

// NewChain returns a in-memory chain that implements the TSDB interface.
func NewChain(path string) *Chain {
	logger.File(fmt.Sprintf("creating new chain at path %s", path), "p")

	return &Chain{
		Path:           path,
		Chain:          []Block{},
		LengthElements: 0,
		Size:           0,
	}
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

	// GetPositionalIndexNormalized accepts the normalized time, searches for the block with that time
	// using jump search, and returns the address of the block having the specified normalized
	// time.
	GetPositionalIndexNormalized(n int64) (int, error)

	// PopPreviousNBlocks pops or removes **n** previous blocks from the chain and returns
	// success status.
	PopPreviousNBlocks(n uint64) (Chain, error)

	// GetChain returns the positional pointer address of the first element of the chain.
	GetChain() *[]Block

	// Commit saves or commits the chain in storage and returns success status.
	Commit() bool
}

// Init initialize Chain properties
func (c *Chain) Init() *Chain {
	c.mux.Lock()
	defer c.mux.Unlock()

	res, e := parse(c.Path)
	if e != nil {
		logger.Terminal(fmt.Sprintf("creating chain at %s", c.Path), "p")
		c.LengthElements = 0
		c.Size = unsafe.Sizeof(c)
		c.Chain = []Block{}
		return c
	}

	raw := loadFromStorage(res)
	c.Chain = *raw
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

// Commit saves or commits the existing chain in the secondary memory.
// Returns the success status
func (c *Chain) Commit() *Chain {
	c.mux.Lock()
	defer c.mux.Unlock()

	logger.File("writing chain of length"+strconv.Itoa(len(c.Chain)), "p")
	bytes := parseToJSON(c.Chain)
	e := saveToHDD(c.Path, bytes)
	if e != nil {
		panic(e)
	}
	return c
}

// GetPositionalIndexNormalized Returns block by searching the chain for the NormalizedTime
func (c *Chain) GetPositionalIndexNormalized(n int64) (int, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.LengthElements = len(c.Chain)
	if c.Chain[c.LengthElements-1].NormalizedTime < n || c.Chain[0].NormalizedTime > n {
		return 0, errors.New("Normalized time not in Chain range")
	}
	l, r, m := 0, c.LengthElements, 0

	for l <= r {
		m = l + (r-l)/2
		if c.Chain[m].NormalizedTime == n {
			return m, nil
		}

		if c.Chain[m].NormalizedTime < n {
			l = m + 1
		} else {
			r = m - 1
		}
	}

	return 0, errors.New("Normalized time not found in chain")
}

// ChainSet is a set of chain that manages the operations related to chains
// on a macro level. These include flushing chains to the storage based on
// regular time intervals or size (to be done). It can delete chains that are
// not active, thus being low on the memory. Scheduling operations on
// time-series values in chain can be done as well with slight customization.
type ChainSet struct {
	FlushDuration time.Duration
	Cmap          map[string]*Chain
}

func NewChainSet(flushDuration time.Duration) *ChainSet {
	return &ChainSet{
		FlushDuration: flushDuration,
	}
}

// Register makes a new property in the Chain map (Cmap) with
// name as Key and Chain address as value respectively. Repeated
// calls with same name will overwrite the chain contents and hence
// not recommended.
func (cs *ChainSet) Register(name string, chainAddress *Chain) {
	cs.Cmap[name] = chainAddress
}

// Get returns the chain corresponding to the passed name.
func (cs *ChainSet) Get(name string) *Chain {
	return cs.Cmap[name]
}

// Run triggers the ChainSet tasks which currently includes
// flushing those chains that have newer blocks only. This is done
// keeping in mind the performance of the system, thus being effective
// on the resources.
func (cs *ChainSet) Run() {}
