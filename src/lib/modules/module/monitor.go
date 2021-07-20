package module

import (
	"context"
	"fmt"

	config "github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/bench-routes/bench-routes/src/lib/log"
	"github.com/bench-routes/bench-routes/src/lib/modules/job"
	"github.com/bench-routes/bench-routes/src/lib/modules/scheduler"
	"github.com/bench-routes/bench-routes/tsdb/file"
)

// Monitor handles monitoring of the endpoints.
type Monitor struct {
	jobs         map[*job.JobInfo]chan<- struct{}
	existingJobs map[string]struct{}
	chainSet     *file.ChainSet
	reload       chan struct{}
	errCh        chan<- error
}

func newMonitorModule(chainSet *file.ChainSet, errCh chan<- error) (*Monitor, error) {
	job := &Monitor{
		jobs:         make(map[*job.JobInfo]chan<- struct{}),
		existingJobs: make(map[string]struct{}),
		reload:       make(chan struct{}),
		chainSet:     chainSet,
		errCh:        errCh,
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
		schdr := scheduler.New(m.jobs)
		go schdr.Run(ctx)
	}
}

// Reload reloads the new config and signals reload channel. It expects the Run() to already be active.
func (m *Monitor) Reload(conf *config.Config) error {
	for i, api := range conf.APIs {
		_, exists := m.existingJobs[api.Name]
		if exists {
			// todo: deletion of jobs that are no longer existing
			log.Info("component", "reload", "msg", "job already exists with name "+api.Name+". Skipping creation.")
			continue
		}
		app, _ := m.chainSet.NewChain(api.Name+"_monitor", api.Protocol+api.Domain+api.Route, false)

		exec, ch, err := job.NewJob("monitor", app, nil, &api)
		if err != nil {
			return fmt.Errorf("error creating # %d job: %s", i, err)
		}
		go exec.Execute(m.errCh)
		m.jobs[exec.Info()] = ch
		m.existingJobs[api.Name] = struct{}{}
	}
	// signaling to reload.
	m.reload <- struct{}{}
	return nil
}

// Stop stops the module.
func (m *Monitor) Stop() {
	close(m.reload)
}
