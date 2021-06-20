package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
	"github.com/bench-routes/bench-routes/src/lib/modules/job"
	"github.com/bench-routes/bench-routes/tsdb/file"
)

func TestScheduler(t *testing.T) {
	conf, err := config.New("./testdata/config.yml")
	if err != nil {
		t.Fatalf("error loading config: %s", err)
	}
	jobs := make(map[*job.JobInfo]chan<- struct{})
	for i, api := range conf.APIs {
		var app file.Appendable
		ch := make(chan struct{})
		// creating the jobs
		exec, err := job.NewJob("machine", app, ch, &api)
		if err != nil {
			fmt.Println(fmt.Errorf("error creating # %d job: %s\n", i, err))
			continue
		}
		// launching the jobs
		go exec.Execute()
		jobs[exec.Info()] = ch
	}
	ctx, cancel := context.WithCancel(context.Background())
	scheduler := NewScheduler(jobs)
	go scheduler.Run(ctx)
	time.Sleep(time.Second * 28)
	cancel()
	time.Sleep(time.Second * 8)
}
