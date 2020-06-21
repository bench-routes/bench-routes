package tsdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/zairza-cetb/bench-routes/src/lib/logger"
)

const (
	// BlockDataSeparator sets a separator for block data value.
	BlockDataSeparator = "|"
	// TSDBFileExtension file extension for db files.
	TSDBFileExtension = ".json"
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
	Path               string
	Name               string
	Route              string
	Chain              []Block
	LengthElements     int
	Size               uintptr
	containsNewBlocks  bool
	inActiveIterations uint32
	mux                sync.Mutex
}

// ChainReadOnly is a read-only structure that contains
// a stream of blocks from the tsdb. This is meant to
// be used by the querier for performing read operations
// over the time-series samples.
type ChainReadOnly struct {
	Path  string
	Chain *[]Block
}

// NewChain returns a in-memory chain that implements the TSDB interface.
func NewChain(path string) *Chain {
	logger.File(fmt.Sprintf("creating new chain at path %s", path), "p")
	return &Chain{
		Name:              filterChainPath(path),
		Path:              path,
		Chain:             []Block{},
		LengthElements:    0,
		Size:              0,
		containsNewBlocks: true,
	}
}

// ReadOnly returns a in-memory chain that implements the TSDB interface.
func ReadOnly(path string) *ChainReadOnly {
	logger.File(fmt.Sprintf("creating new chain at path %s", path), "p")
	var blockStream []Block

	return &ChainReadOnly{
		Path:  path,
		Chain: &blockStream,
	}
}

func filterChainPath(name string) string {
	name = strings.ReplaceAll(name, ".", "___")
	return strings.ReplaceAll(name, "/", "_")
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

	// GetPositionalIndexNormalized accepts the normalized time, searches for the block with that time
	// using jump search, and returns the address of the block having the specified normalized
	// time.
	GetPositionalIndexNormalized(n int64) (int, error)

	// PopPreviousNBlocks pops or removes **n** previous blocks from the chain and returns
	// success status.
	PopPreviousNBlocks(n uint64) (Chain, error)

	// GetChain returns the positional pointer address of the first element of the chain.
	GetChain() *[]Block
}

// Init initialize Chain properties
func (c *Chain) Init() *Chain {
	c.mux.Lock()
	defer c.mux.Unlock()

	res, e := parse(c.Path)
	if e != nil {
		logger.Terminal(fmt.Sprintf("creating in-memory chain: %s", c.Name), "p")
		c.LengthElements = 0
		c.Size = unsafe.Sizeof(c)
		c.Chain = []Block{}
		c.commit()
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
	c.containsNewBlocks = true
	if c.inActiveIterations != 0 {
		c.inActiveIterations = 0
	}
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
func (c *Chain) commit() *Chain {
	logger.File("writing chain of length"+strconv.Itoa(len(c.Chain)), "p")
	bytes := parseToJSON(c.Chain)
	e := saveToHDD(c.Path, bytes)
	if e != nil {
		panic(e)
	}
	c.containsNewBlocks = false
	return c
}

// GetPositionalIndexNormalized Returns block by searching the chain for the NormalizedTime
func (c *Chain) GetPositionalIndexNormalized(n int64) (int, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.LengthElements = len(c.Chain)
	if c.Chain[c.LengthElements-1].NormalizedTime < n || c.Chain[0].NormalizedTime > n {
		return 0, errors.New("normalized time not in Chain range")
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

	return 0, errors.New("normalized time not found in chain")
}

// VerifyChainPathExists verifies the existence of chain in the tsdb directory.
func VerifyChainPathExists(chainPath string) bool {
	_, err := os.Stat(chainPath)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

const (
	// FlushAsTime for flushing in regular intervals of seconds.
	FlushAsTime = 0
	// FlushAsSpace for flushing in regular intervals of space/bytes.
	FlushAsSpace = 1
	// inActiveIterationsLimit is the limit after which the chain is deleted
	// from the chain set in order to free up the memory from inactive
	// chains.
	// inActiveIterationsLimit = 5
)

// ChainSet is a set of chain that manages the operations related to chains
// on a macro level. These include flushing chains to the storage based on
// regular time intervals or size (to be done). It can delete chains that are
// not active, thus being low on the memory. Scheduling operations on
// time-series values in chain can be done as well with slight customization.
type ChainSet struct {
	FlushDuration time.Duration
	flushType     int
	Cmap          map[string]*Chain
	cancel        chan interface{}
	mux           sync.RWMutex
}

// NewChainSet returns a new ChainSet for managing chains during runtime.
func NewChainSet(flushType int, flushDuration time.Duration) *ChainSet {
	return &ChainSet{
		FlushDuration: flushDuration,
		flushType:     flushType,
		Cmap:          make(map[string]*Chain),
		cancel:        make(chan interface{}),
	}
}

// Append currently not supported.
// Appends the block into the chain name passed. The new block is added
// only in the memory. Commit is done by the chain scheduler and only after
// commit, the changes appear in the secondary storage.
func (cs *ChainSet) Append(name string, block Block) *Chain {
	cs.Cmap[name].Append(block)
	return cs.Cmap[name]
}

// Cancel cancels or stops the execution of chain scheduler.
func (cs *ChainSet) Cancel() {
	cs.cancel <- ""
}

// Get returns the chain corresponding to the passed name. It returns
// false if the chain is not found in the Cmap. This can be the case if the
// chain has been deleted by the Run() in order to save the memory resources.
func (cs *ChainSet) Get(name string) (*Chain, bool) {
	if _, ok := cs.Cmap[name]; !ok {
		return nil, false
	}
	return cs.Cmap[name], true
}

// Register makes a new property in the Chain map (Cmap) with
// name as Key and Chain address as value respectively. Repeated
// calls with same name will overwrite the chain contents and hence
// not recommended.
func (cs *ChainSet) Register(name string, chainAddress *Chain) {
	cs.mux.Lock()
	defer cs.mux.Unlock()
	cs.Cmap[name] = chainAddress
}

// Run is a chain scheduler that triggers the ChainSet tasks which currently includes
// flushing those chains that have newer blocks only. This is done
// keeping in mind the performance of the system, thus being effective
// on the resources.
func (cs *ChainSet) Run() {
	switch cs.flushType {
	case FlushAsTime:
		go func() {
			for {
				select {
				case <-cs.cancel:
					return
				default:
				}
				for _, chain := range cs.Cmap {
					if chain.containsNewBlocks {
						chain.commit()
					} else {
						// TODO: delete inactive chains and add them back to Cmap when active.
						chain.inActiveIterations++
					}
				}
				time.Sleep(cs.FlushDuration)
			}
		}()
	case FlushAsSpace:
		// TODO: Support for flushing when the chain content exceeds
		// the limit of bytes.
		return
	}
}

// BlockStream returns the address of stream (or list)
// of blocks from the in-memory chain.
func (c *ChainReadOnly) BlockStream() *[]Block {
	return c.Chain
}

// Refresh loads/reloads the chain from the secondary storage
// which contains the latest samples/blocks.
func (c *ChainReadOnly) Refresh() *ChainReadOnly {
	response, e := parse(c.Path)
	if e != nil {
		logger.Terminal(fmt.Sprintf("error reading the chain: %s", c.Path), "f")
	}

	c.Chain = loadFromStorage(response)
	return c
}

func parse(path string) (*string, error) {
	res, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("file not existed. create a new chain")
	}
	str := string(res)
	return &str, nil
}

// parseToJSON converts the chain into Marshallable JSON
func parseToJSON(a []Block) (j []byte) {
	j, e := json.Marshal(a)
	if e != nil {
		panic(e)
	}
	return
}

func loadFromStorage(raw *string) *[]Block {
	var inst []Block
	b := []byte(*raw)
	e := json.Unmarshal(b, &inst)
	if e != nil {
		panic(e)
	}
	return &inst
}

func checkAndCreatePath(path string) {
	array := strings.Split(path, "/")
	array = array[:len(array)-1]
	path = strings.Join(array, "/")
	_, err := os.Stat(path)
	if err != nil {
		e := os.MkdirAll(path, os.ModePerm)
		if e != nil {
			panic(e)
		}
	}
}

func saveToHDD(path string, data []byte) error {
	checkAndCreatePath(path)
	e := ioutil.WriteFile(path, data, 0755)
	if e != nil {
		return e
	}
	return nil
}

// GetTimeStampCalc returns the timestamp
func GetTimeStampCalc() string {
	t := time.Now()

	return s(t.Year()) + "|" + s(t.Month()) + "|" + s(t.Day()) + "|" + s(t.Hour()) + "|" +
		s(t.Minute()) + "|" + s(t.Second())
}

// GetNormalizedTimeCalc returns the UnixNano time as normalized time.
func GetNormalizedTimeCalc() int64 {
	return time.Now().UnixNano()
}

func s(st interface{}) string {
	return fmt.Sprintf("%v", st)
}
