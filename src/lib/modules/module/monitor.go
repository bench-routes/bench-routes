package module

import (
	"context"
	"fmt"
	"sync"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
	"github.com/bench-routes/bench-routes/src/lib/modules/job"
	"github.com/bench-routes/bench-routes/src/lib/modules/scheduler"
	"github.com/bench-routes/bench-routes/tsdb/file"
)

// Monitor handles monitoring of the endpoints.
type Monitor struct {
	mux    sync.RWMutex
	jobs   map[*job.JobInfo]chan<- struct{}
	reload chan struct{}
}

func newMonitorModule() (*Monitor, error) {
	job := &Monitor{
		jobs:   make(map[*job.JobInfo]chan<- struct{}),
		reload: make(chan struct{}),
	}
	return job, nil
}

// Run listens for reload signal and runs a new scheduler.
func (m *Monitor) Run() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	cancelCurrentScheduler := func(cancel context.CancelFunc) {
		if cancel != nil {
			cancel()
		}
	}
	for {
		_, open := <-m.reload
		if !open {
			cancelCurrentScheduler(cancel)
			break
		}

		// canceling scheduler if already present.
		cancelCurrentScheduler(cancel)
		ctx, cancel = context.WithCancel(context.Background())
		scheduler := scheduler.New(m.jobs)
		go scheduler.Run(ctx)
	}
}

// Reload reloads the new config and signals reload channel.
func (m *Monitor) Reload(conf *config.Config, errCh chan<- error) {
	jobs := make(map[*job.JobInfo]chan<- struct{})
	set := file.NewChainSet(0, time.Second*10)
	set.Run()
	for i, api := range conf.APIs {
		app, _ := set.NewChain(api.Name, api.Domain+api.Route, false)

		// creating the jobs.
		exec, ch, err := job.NewJob("monitor", app, &api)
		if err != nil {
			errCh <- fmt.Errorf("error creating # %d job: %s", i, err)
			continue
		}
		// launching the jobs.
		go exec.Execute()
		jobs[exec.Info()] = ch
	}
	m.mux.Lock()
	m.jobs = jobs
	m.mux.Unlock()
	// signaling to reload.
	m.reload <- struct{}{}
	errCh <- nil
}

// Stop stops the module.
func (m *Monitor) Stop() {
	close(m.reload)
}
