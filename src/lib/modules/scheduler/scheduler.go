package scheduler

import (
	"context"
	"time"

	"github.com/bench-routes/bench-routes/src/lib/modules/job"
)

type Scheduler struct {
	scanFrequency time.Duration
	timeline      map[*job.JobInfo]chan<- struct{}
}

func NewScheduler(jobs map[*job.JobInfo]chan<- struct{}) *Scheduler {
	scheduler := &Scheduler{
		scanFrequency: time.Second,
		timeline:      jobs,
	}
	return scheduler
}

func (s *Scheduler) Run(ctx context.Context) {
	for{
		for info,ch := range s.timeline{
			
		}
		time.Sleep(time.Second)
	}
}
