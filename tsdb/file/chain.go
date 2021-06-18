package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/prometheus/common/log"
)

type Appendable interface {
	// Appends allows to append a block into the chain.
	Append(b Block)
}

type ChainUtils interface {
	// getPath returns the path of the chain.
	Path() string
	// ForceCommit commits the chain irrespective of chain-set's runner.
	ForceCommit()
	// Stream returns the underlying block stream.
	Stream() []Block
}

// chain contains Blocks arranged as a chain
type chain struct {
	mux                sync.RWMutex
	path               string
	name               string
	url                string
	stream             []Block
	containsNewBlocks  bool
	inActiveIterations uint32
}

const storagePrefix = "storage"

// newChain returns a in-memory chain that implements the TSDB interface.
func newChain(name, url string, useTestDir bool) *chain {
	c := &chain{
		name:              name,
		path:              fmt.Sprintf("%s/%s.json", storagePrefix, name),
		url:               url,
		stream:            []Block{NewBlock("null", "")},
		containsNewBlocks: false,
	}
	if useTestDir {
		c.path = fmt.Sprintf("testdata/%s.json", name)
	}
	return c
}

// init initialize Chain properties
func (c *chain) init() {
	data, err := parse(c.path)
	if err != nil {
		log.Infof("creating in-memory chain '%s' with path at '%s'\n", c.name, c.path)
		if err := saveToHDD(c.path, parseToJSON(c.stream)); err != nil {
			panic(err)
		}
		return
	}
	c.stream = loadBlocks(data)
}

func parse(path string) ([]byte, error) {
	res, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("file not existed. create a new chain")
	}
	return res, nil
}

func (c *chain) Path() string {
	return c.path
}

// Append function appends the new block in the chain.
func (c *chain) Append(b Block) {
	c.mux.Lock()
	defer c.mux.Unlock()

	// if c.stream[0].Type == "null" {
	// 	c.stream = c.stream[:0]
	// }
	c.stream = append(c.stream, b)
	c.containsNewBlocks = true
	c.inActiveIterations = 0
}

func (c *chain) ForceCommit() {
	c.commit()
}

// commit saves or commits the existing chain in the secondary memory. It is expected to be called by chainset's run.
func (c *chain) commit() {
	c.mux.Lock()
	defer c.mux.Unlock()

	pathPointer, err := parse(c.path)
	if err != nil {
		panic(err)
	}
	existingBlocks := loadBlocks(pathPointer)
	mergedBlocks := mergeBlocksSlice(existingBlocks, c.stream)
	bytes := parseToJSON(mergedBlocks)
	if err := saveToHDD(c.path, bytes); err != nil {
		panic(err)
	}
	c.stream = make([]Block, 0)
	c.containsNewBlocks = false
}

func (c *chain) Stream() []Block {
	return c.stream
}

// VerifyChainPathExists verifies the existence of chain in the tsdb directory.
func VerifyChainPathExists(chainPath string) bool {
	_, err := os.Stat(chainPath)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func loadBlocks(raw []byte) []Block {
	stream := make([]Block, 0)
	e := json.Unmarshal(raw, &stream)
	if e != nil {
		panic(e)
	}
	return stream
}

func saveToHDD(path string, data []byte) error {
	checkAndCreatePath(path)
	e := ioutil.WriteFile(path, data, 0755)
	if e != nil {
		return e
	}
	return nil
}
