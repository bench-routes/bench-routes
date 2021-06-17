package scheduler

import (
	"time"

	"github.com/bench-routes/bench-routes/src/lib/modules/job"
)

type scheduler struct {
	scanFrequency time.Duration
	timeline      map[*job.JobInfo]chan<- struct{}
}

func NewScheduler() *scheduler {
	scheduler := &scheduler{
		scanFrequency: time.Second,
		timeline:      make(map[*job.JobInfo]chan<- struct{}),
	}
	return scheduler
}

func (s *scheduler) Run() {

}
