package main

import (
	"log"
	"github.com/zairza-cetb/bench-routes/src/collector/process"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

const (
	storagePath = "collector-store/"
)

var (
	// COMMAND: chain
	processChains = make(map[string]tsdb.Chain)
)

func main() {
	log.Printf("Enabling bench-routes process collections")

	buffer := process.NewProcessReader()
	buffer.UpdateCurrentProcesses()

	
}
