package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zairza-cetb/bench-routes/src/service"
	"golang.org/x/sync/errgroup"
)

func main() {
	s := &service.Service{
		Port: ":9090",
	}
	if len(os.Args) > 1 {
		s.Port = ":" + os.Args[1]
	}
	ctx, cancelFunc := context.WithCancel(context.Background())

	g := errgroup.Group{}
	g.Go(func() error {
		return s.Runner(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Println("program terminated")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(c)
		cancelFunc()
	}()

	<-c
	cancelFunc()
}
