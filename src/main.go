package main

import (
	"github.com/zairza-cetb/bench-routes/src/service"
	"os"
)

func main() {

	port := ":9090"
	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}

	service.Runner(port)

}
