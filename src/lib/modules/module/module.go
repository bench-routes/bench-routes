package module

import (
	"context"
	"fmt"
	"sync"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
	"github.com/bench-routes/bench-routes/src/lib/modules/job"
	"github.com/bench-routes/bench-routes/src/lib/modules/scheduler"
	"github.com/bench-routes/bench-routes/tsdb/file"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

type Runnable interface {
	Run()
	Reload(*config.Config)
	Stop()
}

type Machine struct {
	mux    sync.RWMutex
	jobs   map[*job.JobInfo]chan<- struct{}
	reload chan struct{}
}

type Monitor struct {
	mux    sync.RWMutex
	jobs   map[*job.JobInfo]chan<- struct{}
	reload chan struct{}
}

func NewModule(typ string) (Runnable, error) {
	switch typ {
	case "machine":
		module, err := newMachineModule()
		if err != nil {
			return nil, fmt.Errorf("error creating module : %w", err)
		}
		return module, nil
	case "monitor":
		module, err := newMonitorModule()
		if err != nil {
			return nil, fmt.Errorf("error creating module : %w", err)
		}
		return module, nil
	default:
		return nil, fmt.Errorf("`typ` provided is invalid")
	}
}

func newMachineModule() (*Machine, error) {
	job := &Machine{
		jobs:   make(map[*job.JobInfo]chan<- struct{}),
		reload: make(chan struct{}),
	}
	return job, nil
}

func newMonitorModule() (*Monitor, error) {
	job := &Monitor{
		jobs:   make(map[*job.JobInfo]chan<- struct{}),
		reload: make(chan struct{}),
	}
	return job, nil
}

func (m *Machine) Run() {
	for {
		_, open := <-m.reload
		if !open {
			if cancel != nil {
				cancel()
			}
			fmt.Println("Stopping Module")
			break
		}
		if cancel != nil {
			cancel()
		}
		ctx, cancel = context.WithCancel(context.Background())
		scheduler := scheduler.NewScheduler(m.jobs)
		go scheduler.Run(ctx)
	}
}

func (m *Machine) Reload(conf *config.Config) {
	jobs := make(map[*job.JobInfo]chan<- struct{})
	for i, api := range conf.APIs {
		var app file.Appendable
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
	// signalling to reload the scheduler
	m.reload <- struct{}{}
}

func (m *Machine) Stop() {
	close(m.reload)
}

func (m *Monitor) Run() {
	for {
		_, open := <-m.reload
		if !open {
			if cancel != nil {
				cancel()
			}
			break
		}
		if cancel != nil {
			cancel()
		}
		ctx, cancel = context.WithCancel(context.Background())
		scheduler := scheduler.NewScheduler(m.jobs)
		go scheduler.Run(ctx)
	}
}

func (m *Monitor) Reload(conf *config.Config) {
	jobs := make(map[*job.JobInfo]chan<- struct{})
	for i, api := range conf.APIs {
		var app file.Appendable
		ch := make(chan struct{})
		// creating the jobs
		exec, err := job.NewJob("monitor", app, ch, &api)
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
	// signalling to reload the scheduler
	m.reload <- struct{}{}
}

func (m *Monitor) Stop() {
	close(m.reload)
}
