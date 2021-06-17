package module

import (
	"fmt"
	"sync"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
	"github.com/bench-routes/bench-routes/src/lib/modules/job"
	"github.com/bench-routes/bench-routes/tsdb/file"
)

type Runnable interface {
	Run()
	Reload(*config.Config)
	Stop()
}

type Machine struct {
	// scheduler
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
	// case "monitor":
	// 	module, err := newMonitorModule(conf)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error creating module : %w", err)
	// 	}
	// 	return module, nil
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

func (m *Machine) Run() {
	for range m.reload {

	}
}

func (m *Machine) Reload(conf *config.Config) {
	jobs := make(map[*job.JobInfo]chan<- struct{})
	for _, api := range conf.APIs {
		var app file.Appendable
		ch := make(chan struct{})
		exec, err := job.NewJob("machine", &app, ch, &api)
		if err != nil {
			fmt.Println(fmt.Errorf("error creating job: %s", err))
		}
		go exec.Execute()
		jobs[exec.Info()] = ch
	}
	m.mux.Lock()
	m.jobs = jobs
	m.mux.Unlock()
	m.reload <- struct{}{}
}

func (m *Machine) Stop() {
	close(m.reload)
}
