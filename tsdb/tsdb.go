package tsdb

import (
	"errors"
	"log"
	"math"
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

// PingTSDB type for PingTSDB
// type PingTSDB struct {
// 	URL   string
// 	Hash  string
// 	Path  string
// 	Chain ChainPing
// }

// FloodPingTSDB type for PingTSDB
type FloodPingTSDB struct {
	URL   string
	Hash  string
	Path  string
	Chain ChainFloodPing
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

// FloodPingType type for storing Ping values in TSDB
type FloodPingType struct {
	Min        float64 `yaml:"min"`
	Mean       float64 `yaml:"mean"`
	Max        float64 `yaml:"max"`
	MDev       float64 `yaml:"mdev"`
	PacketLoss float64 `yaml:"packetloss"`
}

// BlockPing block for ping case
type BlockPing struct {
	PrevBlock      *BlockPing
	NextBlock      *BlockPing
	Datapoint      PingType
	NormalizedTime int64
	Timestamp      time.Time
}

// BlockFloodPing block for ping case
type BlockFloodPing struct {
	PrevBlock      *BlockFloodPing
	NextBlock      *BlockFloodPing
	Datapoint      FloodPingType
	NormalizedTime int64
	Timestamp      time.Time
}

// BlockPingJSON block for ping case
type BlockPingJSON struct {
	Datapoint      PingType  `json:"datapoint"`
	NormalizedTime int64     `json:"normalizedTime"`
	Timestamp      time.Time `json:"timestamp"`
}

// BlockFloodPingJSON block for flood-ping case
type BlockFloodPingJSON struct {
	Datapoint      FloodPingType `json:"datapoint"`
	NormalizedTime int64         `json:"normalizedTime"`
	Timestamp      time.Time     `json:"timestamp"`
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
	mux            sync.Mutex
}

// ChainPing forms chain for Ping values
type ChainPing struct {
	Path           string
	Chain          []BlockPing
	LengthElements int
	Size           uintptr
	mux            sync.Mutex
}

// ChainFloodPing forms chain for Ping values
type ChainFloodPing struct {
	Path           string
	Chain          []BlockFloodPing
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
	// For Ping module
	InitPing() (*[]BlockPing, ChainPing)
	// For Flood Ping module
	InitFloodPing() (*[]BlockFloodPing, ChainFloodPing)

	// Append appends a new tsdb block passed as params to the most recent location (or
	// the last location) of the chain. Returns success status.
	Append(b Block) bool
	// For Ping module
	AppendPing(b BlockPing) bool
	// For Flood Ping module
	AppendFloodPing(b BlockFloodPing) bool

	// GetPositionalPointerNormalized accepts the normalized time, searches for the block with that time
	// using jump search, and returns the address of the block having the specified normalized
	// time.
	GetPositionalPointerNormalized(n int64) (*Block, error)
	// For Ping module
	GetPositionalPointerNormalizedPing(n int64) (*BlockPing, error)
	// For Flood Ping module
	GetPositionalPointerNormalizedFloodPing(n int64) (*BlockFloodPing, error)

	// PopPreviousNBlocks pops or removes **n** previous blocks from the chain and returns
	// success status.
	PopPreviousNBlocks(n uint64) (Chain, error)
	// For Ping module
	PopPreviousNBlocksPing(n uint64) (ChainPing, error)
	// For Flood Ping module
	PopPreviousNBlocksFloodPing(n uint64) (ChainFloodPing, error)

	// GetChain returns the positional pointer address of the first element of the chain.
	GetChain() *[]Block
	// For Ping module
	GetChainPing() *[]BlockPing
	// For Flood Ping module
	GetChainFloodPing() *[]BlockFloodPing

	// Save saves or commits the chain in storage and returns success status.
	Save() bool
	// For Ping module
	SavePing() bool
	// For Flood Ping module
	SaveFloodPing() bool
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

// InitPing initialize Ping Chain properties
func (c *ChainPing) InitPing() *ChainPing {
	c.mux.Lock()
	defer c.mux.Unlock()
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

// InitFloodPing initialize Ping Chain properties
func (c *ChainFloodPing) InitFloodPing() *ChainFloodPing {
	c.mux.Lock()
	defer c.mux.Unlock()
	res, e := parse(c.Path)
	if e != nil {
		log.Printf("chain not found at %s. creating one ...", c.Path)
		c.LengthElements = 0
		c.Size = unsafe.Sizeof(c)
		c.Chain = []BlockFloodPing{}
		return c
	}

	raw := loadFromStorageFloodPing(res)
	c.Chain = *formLinkedChainFromRawBlockFloodPing(raw)
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

func formLinkedChainFromRawBlockFloodPing(a *[]BlockFloodPingJSON) *[]BlockFloodPing {
	r := *a
	l := len(r)
	arr := []BlockFloodPing{}
	for i := 0; i < l; i++ {
		inst := BlockFloodPing{
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
func (c *Chain) Append(b Block) *Chain {

	c.mux.Lock()
	defer c.mux.Unlock()
	c.Chain = append(c.Chain, b)
	c.Size = unsafe.Sizeof(c)
	c.LengthElements = len(c.Chain)
	l := c.LengthElements
	if l != 1 {
		c.Chain[l-2].NextBlock = &c.Chain[l-1]
		c.Chain[l-1].PrevBlock = &c.Chain[l-2]
		c.Chain[l-1].NextBlock = nil
		// update normalised time
		c.Chain[l-1].NormalizedTime = c.Chain[l-2].NormalizedTime + 1
	}
	return c

}

// AppendPing function appends the new ping block in the chain
func (c *ChainPing) AppendPing(b BlockPing) *ChainPing {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Chain = append(c.Chain, b)
	c.Size = unsafe.Sizeof(c)
	c.LengthElements = len(c.Chain)
	l := c.LengthElements
	if l != 1 {
		c.Chain[l-2].NextBlock = &c.Chain[l-1]
		c.Chain[l-1].PrevBlock = &c.Chain[l-2]
		c.Chain[l-1].NextBlock = nil
		// update normalised time
		c.Chain[l-1].NormalizedTime = c.Chain[l-2].NormalizedTime + 1
	}
	return c

}

// AppendFloodPing function appends the new flood ping block in the chain
func (c *ChainFloodPing) AppendFloodPing(b BlockFloodPing) *ChainFloodPing {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Chain = append(c.Chain, b)
	c.Size = unsafe.Sizeof(c)
	c.LengthElements = len(c.Chain)
	l := c.LengthElements
	if l != 1 {
		c.Chain[l-2].NextBlock = &c.Chain[l-1]
		c.Chain[l-1].PrevBlock = &c.Chain[l-2]
		c.Chain[l-1].NextBlock = nil
		// update normalised time
		c.Chain[l-1].NormalizedTime = c.Chain[l-2].NormalizedTime + 1
	}
	return c

}

// PopPreviousNBlocks pops last n elements from chain.
func (c *Chain) PopPreviousNBlocks(n int) (*Chain, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
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
func (c *ChainPing) PopPreviousNBlocksPing(n int) (*ChainPing, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
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

// PopPreviousNBlocksFloodPing pops last n elements from flood ping chain.
func (c *ChainFloodPing) PopPreviousNBlocksFloodPing(n int) (*ChainFloodPing, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
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

// SavePing saves or commits the existing chain in the secondary memory.
// Returns the success status
func (c *ChainPing) SavePing() bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	log.Printf("writing Ping chain of length %d", len(c.Chain))
	bytes := parser.ParseToJSONPing(c.Chain)
	e := saveToHDD(c.Path, bytes)
	if e != nil {
		panic(e)
	}
	return true
}

// SaveFloodPing saves or commits the existing chain in the secondary memory.
// Returns the success status
func (c *ChainFloodPing) SaveFloodPing() bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	log.Printf("writing Ping chain of length %d", len(c.Chain))
	bytes := parser.ParseToJSONFloodPing(c.Chain)
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
func (c *ChainPing) GetPositionalPointerNormalizedPing(n int64) (*BlockPing, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
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

// GetPositionalPointerNormalizedFloodPing Returns block by searching the chain for the NormalizedTime
func (c *ChainFloodPing) GetPositionalPointerNormalizedFloodPing(n int64) (*BlockFloodPing, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
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
