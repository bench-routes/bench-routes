package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/bench-routes/bench-routes/src/lib/modules/job"
)

type scheduler struct {
	scanFrequency time.Duration
	timeline      map[*job.JobInfo]chan<- struct{}
}

// New returns a scheduler with new timeline.
func New(jobs map[*job.JobInfo]chan<- struct{}) *scheduler {
	s := &scheduler{
		scanFrequency: time.Second * 1,
		timeline:      jobs,
	}
	return s
}

// Run runs the scheduler with a ticker of one second.
func (s *scheduler) Run(ctx context.Context) {
	d := s.scanFrequency
	for {
		select {
		case <-time.After(d):
			for info, ch := range s.timeline {
				// We do not stop the scheduler even if we get a cancel request while scheduling
				// the jobs. The cancel should happen only when all the jobs have been served
				// and then come out (or get cancelled) in the next iteration. This keeps the system
				// with a deterministic behaviour.
				if info.Every <= time.Since(info.ReadTime()) {
					ch <- struct{}{}
				}
			}
		case <-ctx.Done():
			log.Println("Stopping Scheduler")
			return
		}
	}
}
