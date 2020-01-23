package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/collector/process"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

const (
	storagePath = "collector-store/"
	defaultWait = time.Second * 10
)

var (
	// COMMAND: chain
	processChains = make(map[string]*tsdb.Chain)
)

func main() {
	log.Printf("Enabling bench-routes process collections")

	buffer := process.NewProcessReader()
	if _, err := buffer.UpdateCurrentProcesses(); err != nil {
		panic(err)
	}

	for _, pdetails := range *buffer.ProcessesDetails {
		processChains[pdetails.FilteredCommand] = tsdb.NewChain(fmt.Sprintf("%s%s.json", storagePath, pdetails.FilteredCommand))
		processChains[pdetails.FilteredCommand].Init().Commit()
	}

	for {
		var wg sync.WaitGroup
		if _, err := buffer.UpdateCurrentProcesses(); err != nil {
			panic(err)
		}

		wg.Add(buffer.TotalRunningProcesses)

		for _, ps := range *buffer.ProcessesDetails {
			go func(ps process.PDetails) {
				b := tsdb.GetNewBlock("ps", ps.Encode())
				processChains[ps.FilteredCommand].Append(*b).Commit()
				wg.Done()
			}(ps)
		}

		wg.Wait()
		time.Sleep(defaultWait)
	}
}
