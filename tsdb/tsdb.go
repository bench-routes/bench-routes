package tsdb

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"
	"unsafe"
)

var (
	parser Parser
	// PingDBNames contains the name of the database corresponding to the unique config url
	PingDBNames = make(map[string]string)
	// GlobalPingChain contains chains of all the pings operating in bench-routes which has to be globally accessed
	GlobalPingChain []ChainPing
)

// PingTSDB type for PingTSDB
type PingTSDB struct {
	URL   string
	Hash  string
	Path  string
	Chain ChainPing
}

// Block use case block for the TSDB chain
type Block struct {
	PrevBlock      *Block
	NextBlock      *Block
	Datapoint      float32
	NormalizedTime int64
	Timestamp      time.Time
}

// PingType type for storing Ping values in TSDB
type PingType struct {
	Min  float64 `yaml:"min"`
	Mean float64 `yaml:"mean"`
	Max  float64 `yaml:"max"`
	MDev float64 `yaml:"mdev"`
}

// BlockPing block for ping case
type BlockPing struct {
	PrevBlock      *BlockPing
	NextBlock      *BlockPing
	Datapoint      PingType
	NormalizedTime int64
	Timestamp      time.Time
}

// BlockPingJSON block for ping case
type BlockPingJSON struct {
	Datapoint      PingType  `json:"datapoint"`
	NormalizedTime int64     `json:"normalizedTime"`
	Timestamp      time.Time `json:"timestamp"`
}

// BlockJSON helps reffer Block as JSON
type BlockJSON struct {
	Datapoint      float32   `json:"datapoint"`
	NormalizedTime int64     `json:"normalizedTime"`
	Timestamp      time.Time `json:"timestamp"`
}

// Chain contains Blocks arranged as a chain
type Chain struct {
	Path           string
	Chain          []Block
	LengthElements int
	Size           uintptr
}

// ChainPing forms chain for Ping values
type ChainPing struct {
	Path           string
	Chain          []BlockPing
	LengthElements int
	Size           uintptr
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
	// For Ping module
	InitPing() (*[]BlockPing, ChainPing)

	// Append appends a new tsdb block passed as params to the most recent location (or
	// the last location) of the chain. Returns success status.
	Append(b Block) bool
	// For Ping module
	AppendPing(b BlockPing) bool

	// GetPositionalPointerNormalized accepts the normalized time, searches for the block with that time
	// using jump search, and returns the address of the block having the specified normalized
	// time.
	GetPositionalPointerNormalized(n int64) (*Block, error)
	// For Ping module
	GetPositionalPointerNormalizedPing(n int64) (*BlockPing, error)

	// PopPreviousNBlocks pops or removes **n** previous blocks from the chain and returns
	// success status.
	PopPreviousNBlocks(n uint64) (Chain, error)
	// For Ping module
	PopPreviousNBlocksPing(n uint64) (ChainPing, error)

	// GetChain returns the positional pointer address of the first element of the chain.
	GetChain() *[]Block
	// For Ping module
	GetChainPing() *[]BlockPing

	// Save saves or commits the chain in storage and returns success status.
	Save() bool
	// For Ping module
	SavePing() bool
}

// Init initialize Chain properties
func (c Chain) Init() Chain {
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

// InitPing initialize Ping Chain properties
func (c ChainPing) InitPing() ChainPing {
	res, e := parse(c.Path)
	if e != nil {
		log.Printf("chain not found at %s. creating one ...", c.Path)
		c.LengthElements = 0
		c.Size = unsafe.Sizeof(c)
		c.Chain = []BlockPing{}
		return c
	}

	raw := loadFromStoragePing(res)
	c.Chain = *formLinkedChainFromRawBlockPing(raw)
	c.LengthElements = len(c.Chain)
	c.Size = unsafe.Sizeof(c)
	return c
}

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

func formLinkedChainFromRawBlockPing(a *[]BlockPingJSON) *[]BlockPing {
	r := *a
	l := len(r)
	arr := []BlockPing{}
	for i := 0; i < l; i++ {
		inst := BlockPing{
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

// Append function appends the new block in the chain
func (c Chain) Append(b Block) Chain {

	c.Chain = append(c.Chain, b)
	c.Size = unsafe.Sizeof(c)
	c.LengthElements = len(c.Chain)
	l := c.LengthElements
	if l != 1 {
		c.Chain[l-2].NextBlock = &c.Chain[l-1]
		c.Chain[l-1].PrevBlock = &c.Chain[l-2]
		c.Chain[l-1].NextBlock = nil
	}
	return c

}

// AppendPing function appends the new ping block in the chain
func (c ChainPing) AppendPing(b BlockPing) ChainPing {

	fmt.Println("path hererrer is ", c.Path)
	c.Chain = append(c.Chain, b)
	c.Size = unsafe.Sizeof(c)
	c.LengthElements = len(c.Chain)
	l := c.LengthElements
	if l != 1 {
		c.Chain[l-2].NextBlock = &c.Chain[l-1]
		c.Chain[l-1].PrevBlock = &c.Chain[l-2]
		c.Chain[l-1].NextBlock = nil
	}
	return c

}

// PopPreviousNBlocks pops last n elements from chain.
func (c Chain) PopPreviousNBlocks(n int) (Chain, error) {
	c.LengthElements = len(c.Chain)
	l := c.LengthElements
	if c.Chain[l-1].NextBlock != nil || c.Chain[0].PrevBlock != nil {
		return c, errors.New("Chain corrupted")
	} else if l < n {
		return c, errors.New("Deletion will cause underflow")
	}
	c.Chain = c.Chain[:len(c.Chain)-n]
	c.Chain[l-n-1].NextBlock = nil
	c.LengthElements = l - n
	c.Size = unsafe.Sizeof(c)
	return c, nil
}

// PopPreviousNBlocksPing pops last n elements from ping chain.
func (c ChainPing) PopPreviousNBlocksPing(n int) (ChainPing, error) {
	c.LengthElements = len(c.Chain)
	l := c.LengthElements
	if c.Chain[l-1].NextBlock != nil || c.Chain[0].PrevBlock != nil {
		return c, errors.New("Chain corrupted")
	} else if l < n {
		return c, errors.New("Deletion will cause underflow")
	}
	c.Chain = c.Chain[:len(c.Chain)-n]
	c.Chain[l-n-1].NextBlock = nil
	c.LengthElements = l - n
	c.Size = unsafe.Sizeof(c)
	return c, nil
}

// Save saves or commits the existing chain in the secondary memory.
// Returns the success status
func (c Chain) Save() bool {
	bytes := parser.ParseToJSON(c.Chain)
	e := saveToHDD(c.Path, bytes)
	if e != nil {
		panic(e)
	}
	return true
}

// SavePing saves or commits the existing chain in the secondary memory.
// Returns the success status
func (c ChainPing) SavePing() bool {
	fmt.Println("SAVING PATH IS ", c.Path)
	bytes := parser.ParseToJSONPing(c.Chain)
	e := saveToHDD(c.Path, bytes)
	if e != nil {
		panic(e)
	}
	return true
}

// GetPositionalPointerNormalized Returns block by searching the chain for the NormalizedTime
func (c Chain) GetPositionalPointerNormalized(n int64) (*Block, error) {
	c.LengthElements = len(c.Chain)
	jumpSize := int(math.Floor(math.Sqrt(float64(c.LengthElements))))
	index := jumpSize - 1

	if c.Chain[c.LengthElements-1].NormalizedTime < n || c.Chain[0].NormalizedTime > n {
		return nil, errors.New("Normalized time not in Chain range")
	}

	for n > c.Chain[index].NormalizedTime && index < c.LengthElements-jumpSize {
		index += jumpSize
	}

	for c.Chain[index].NormalizedTime < n {
		index++
	}

	for c.Chain[index].NormalizedTime > n {
		index--
	}

	if c.Chain[index].NormalizedTime != n {
		return nil, errors.New("Normalized time not found in chain")
	}

	return &c.Chain[index], nil
}

// GetPositionalPointerNormalizedPing Returns block by searching the chain for the NormalizedTime
func (c ChainPing) GetPositionalPointerNormalizedPing(n int64) (*BlockPing, error) {
	c.LengthElements = len(c.Chain)
	jumpSize := int(math.Floor(math.Sqrt(float64(c.LengthElements))))
	index := jumpSize - 1

	if c.Chain[c.LengthElements-1].NormalizedTime < n || c.Chain[0].NormalizedTime > n {
		return nil, errors.New("Normalized time not in Chain range")
	}

	for n > c.Chain[index].NormalizedTime && index < c.LengthElements-jumpSize {
		index += jumpSize
	}

	for c.Chain[index].NormalizedTime < n {
		index++
	}

	for c.Chain[index].NormalizedTime > n {
		index--
	}

	if c.Chain[index].NormalizedTime != n {
		return nil, errors.New("Normalized time not found in chain")
	}

	return &c.Chain[index], nil
}
