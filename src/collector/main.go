package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/collector/process"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

const (
	storagePath = "collector-store/"
	defaultWait = time.Second * 1
)

var (
	// key: COMMAND; value: chain
	processChains = make(map[string]*tsdb.Chain)
)

func main() {
	log.Printf("Enabling bench-routes process collections")

	buffer := process.NewProcessReader()
	if _, err := buffer.UpdateCurrentProcesses(); err != nil {
		panic(err)
	}

	for _, pdetails := range *buffer.ProcessesDetails {
		p := fmt.Sprintf("%s%s.json", storagePath, pdetails.FilteredCommand)
		assignChaintoMap(&processChains, pdetails.FilteredCommand, p)
	}

	for {
		var wg sync.WaitGroup
		if _, err := buffer.UpdateCurrentProcesses(); err != nil {
			panic(err)
		}

		wg.Add(buffer.TotalRunningProcesses)

		for _, ps := range *buffer.ProcessesDetails {
			go func(ps process.PDetails) {
				if processChains[ps.FilteredCommand] == nil {
					p := fmt.Sprintf("%s%s.json", storagePath, ps.FilteredCommand)
					assignChaintoMap(&processChains, ps.FilteredCommand, p)
				}
				b := *tsdb.GetNewBlock("ps", ps.Encode())
				processChains[ps.FilteredCommand].Append(b).Commit()
				wg.Done()
			}(ps)
		}

		runtime.GC()

		wg.Wait()
		time.Sleep(defaultWait)
	}
}

func assignChaintoMap(c *map[string]*tsdb.Chain, n, path string) {
	(*c)[n] = tsdb.NewChain(path)
	(*c)[n].Init().Commit()
}
