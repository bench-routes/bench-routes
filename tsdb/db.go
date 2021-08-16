package file

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var Chains *ChainSet

func init() {
	Chains = NewChainSet(FlushAsTime, time.Second*10)
	Chains.Run()
}

const (
	// BlockDataSeparator sets a separator for block data value.
	BlockDataSeparator = "|"
	// FileExtension file extension for db files.
	FileExtension = ".json"
)

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
	Cmap          map[string]*chain
	cancel        chan interface{}
	mux           sync.RWMutex
}

// NewChainSet returns a new ChainSet for managing chains during runtime.
func NewChainSet(flushType int, flushDuration time.Duration) *ChainSet {
	return &ChainSet{
		FlushDuration: flushDuration,
		flushType:     flushType,
		Cmap:          make(map[string]*chain),
		cancel:        make(chan interface{}),
	}
}

// Cancel cancels or stops the execution of chain scheduler.
func (cs *ChainSet) Cancel() {
	cs.cancel <- ""
}

// Get returns the chain corresponding to the passed name. It returns
// false if the chain is not found in the Cmap. This can be the case if the
// chain has been deleted by the Run() in order to save the memory resources.
func (cs *ChainSet) Get(name string) (*chain, bool) {
	if _, ok := cs.Cmap[name]; !ok {
		return nil, false
	}
	return cs.Cmap[name], true
}

// NewChain returns a new in-memory chain after registering in chain-set.
func (cs *ChainSet) NewChain(name, url string, useTestDir bool) (Appendable, ChainUtils) {
	c := newChain(name, url, useTestDir)
	c.init()
	cs.register(name, c)
	return c, c
}

// DeleteChain removes the chain.
func (cs *ChainSet) DeleteChain(name string) error {
	_, exists := cs.Cmap[name]
	if !exists {
		return fmt.Errorf("chain '%s' does not exists", name)
	}
	cs.Cmap[name] = nil
	delete(cs.Cmap, name)
	return nil
}

// register makes a new property in the Chain map (Cmap) with
// name as Key and Chain address as value respectively. Repeated
// calls with same name will overwrite the chain contents and hence
// not recommended.
func (cs *ChainSet) register(name string, c *chain) {
	cs.mux.Lock()
	defer cs.mux.Unlock()
	cs.Cmap[name] = c
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
				cs.mux.Lock()
				for _, chain := range cs.Cmap {
					chain.mux.Lock()
					if chain.containsNewBlocks {
						chain.commit()
					} else {
						// TODO: delete inactive chains and add them back to Cmap when active.
						chain.inActiveIterations++
					}
					chain.mux.Unlock()
				}
				cs.mux.Unlock()
				runtime.GC()
				time.Sleep(cs.FlushDuration)
			}
		}()
	case FlushAsSpace:
		// TODO: Support for flushing when the chain content exceeds
		// the limit of bytes.
		return
	}
}

// parseToJSON converts the chain into Marshallable JSON.
func parseToJSON(a []Block) (j []byte) {
	j, e := json.Marshal(a)
	if e != nil {
		panic(e)
	}
	return
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

// GetTimeStampCalc returns the timestamp
func GetTimeStampCalc() string {
	t := time.Now()

	return s(t.Year()) + "|" + s(t.Month()) + "|" + s(t.Day()) + "|" + s(t.Hour()) + "|" +
		s(t.Minute()) + "|" + s(t.Second())
}

// CalcTimeStamp returns the timestamp `sub` seconds previous
func CalcTimeStamp(sub int) string {
	t := time.Now().Add(time.Second * time.Duration(-sub))

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
