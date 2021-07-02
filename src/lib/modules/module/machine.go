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

// Machine handles scraping ping and jitter of the endpoints.
type Machine struct {
	mux    sync.RWMutex
	jobs   map[*job.JobInfo]chan<- struct{}
	reload chan struct{}
}

func newMachineModule() (*Machine, error) {
	job := &Machine{
		jobs:   make(map[*job.JobInfo]chan<- struct{}),
		reload: make(chan struct{}),
	}
	return job, nil
}

// Run listens for reload signal and runs a new scheduler.
func (m *Machine) Run() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	for {
		_, open := <-m.reload
		if !open {
			if cancel != nil {
				cancel()
			}
			break
		}

		// canceling scheduler if already present
		if cancel != nil {
			cancel()
		}
		ctx, cancel = context.WithCancel(context.Background())
		scheduler := scheduler.NewScheduler(m.jobs)
		go scheduler.Run(ctx)
	}
}

// Reload reloads the new config and signals reload channel.
func (m *Machine) Reload(conf *config.Config) {
	jobs := make(map[*job.JobInfo]chan<- struct{})
	set := file.NewChainSet(0, time.Second*10)
	set.Run()
	for i, api := range conf.APIs {
		app, _ := set.NewChain(api.Name, api.Domain+api.Route, false)
		ch := make(chan struct{})
		// creating the jobs
		exec, err := job.NewJob("machine", app, ch, &api)
		if err != nil {
			fmt.Println(fmt.Errorf("error creating # %d job: %s", i, err))
			continue
		}
		// launching the jobs
		go exec.Execute()
		jobs[exec.Info()] = ch
	}
	m.mux.Lock()
	m.jobs = jobs
	m.mux.Unlock()
	// signaling to reload the scheduler
	m.reload <- struct{}{}
}

// Stop stops the module.
func (m *Machine) Stop() {
	close(m.reload)
}