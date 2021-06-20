package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/bench-routes/bench-routes/src/lib/modules/job"
)

type Scheduler struct {
	scanFrequency time.Duration
	timeline      map[*job.JobInfo]chan<- struct{}
}

func NewScheduler(jobs map[*job.JobInfo]chan<- struct{}) *Scheduler {
	scheduler := &Scheduler{
		scanFrequency: time.Second*1 + time.Second/2,
		timeline:      jobs,
	}
	return scheduler
}

func (s *Scheduler) Run(ctx context.Context) {
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
