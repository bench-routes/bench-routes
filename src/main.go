package main

import (
	"github.com/zairza-cetb/bench-routes/src/service"
	"os"
)

var (
	port = ":9090"
)

func main() {

	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}

	service.Runner(port)

}
