package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/bench-routes/bench-routes/src/lib/modules/job"
)

type scheduler struct {
	scanFrequency time.Duration
	timeline      map[*job.JobInfo]chan<- struct{}
}

// NewScheduler returns a scheduler with new timeline
func NewScheduler(jobs map[*job.JobInfo]chan<- struct{}) *scheduler {
	s := &scheduler{
		scanFrequency: time.Second * 1,
		timeline:      jobs,
	}
	return s
}

// Run runs the scheduler with a ticker of one second
func (s *scheduler) Run(ctx context.Context) {
	var d time.Duration = s.scanFrequency
	for {
		select {
		case <-time.After(d):
			fmt.Printf("Rechecking APIs : %s\n", time.Now().Format("15:04:05"))
			for info, ch := range s.timeline {
				if info.Every <= time.Since(info.ReadTime()) {
					fmt.Printf("Execute %s : %s\n", info.Name, info.ReadTime().Format("15:04:05"))
					ch <- struct{}{}
				}
			}
		case <-ctx.Done():
			fmt.Println("Stopping Scheduler")
			return
		}
	}
}
