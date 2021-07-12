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

// Machine handles scraping ping and jitter of the endpoints.
type Machine struct {
	jobs         map[*job.JobInfo]chan<- struct{}
	existingJobs map[string]struct{}
	reload       chan struct{}
	chainSet     *file.ChainSet
	errCh        chan<- error
}

func newMachineModule(chainSet *file.ChainSet, errCh chan<- error) (*Machine, error) {
	job := &Machine{
		jobs:         make(map[*job.JobInfo]chan<- struct{}),
		existingJobs: make(map[string]struct{}),
		reload:       make(chan struct{}),
		chainSet:     chainSet,
		errCh:        errCh,
	}
	return job, nil
}

// Run listens for reload signal and runs a new scheduler.
func (m *Machine) Run() {
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
		schd := scheduler.New(m.jobs)
		go schd.Run(ctx)
	}
}

// Reload reloads the new config and signals reload channel.
func (m *Machine) Reload(conf *config.Config) error {
	for i, api := range conf.APIs {
		_, exists := m.existingJobs[api.Name]
		if exists {
			// todo: deletion of jobs that are no longer existing
			log.Info("component", "reload", "msg", "job already exists with name "+api.Name+". Skipping creation.")
			continue
		}
		app, _ := m.chainSet.NewChain(api.Name+"_machine", api.Domain+api.Route, false)
		exec, ch, err := job.NewJob("machine", app, &api)
		if err != nil {
			return fmt.Errorf("error creating # %d job: %s", i, err)
		}
		// launching the jobs
		go exec.Execute(m.errCh)
		m.jobs[exec.Info()] = ch
		m.existingJobs[api.Name] = struct{}{}
	}
	// signaling to reload the scheduler.
	m.reload <- struct{}{}
	return nil
}

// Stop stops the module.
func (m *Machine) Stop() {
	close(m.reload)
}
