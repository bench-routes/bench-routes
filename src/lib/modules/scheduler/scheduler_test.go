package scheduler

import (
	"context"
	"os"
	"testing"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/bench-routes/bench-routes/src/lib/modules/job"
	"github.com/bench-routes/bench-routes/tsdb/file"
	"github.com/stretchr/testify/require"
)

func TestScheduler(t *testing.T) {
	conf, err := config.New("./testdata/config.yml")
	if err != nil {
		require.FailNow(t, "error loading config: %w\n", err)
	}
	jobs := make(map[*job.JobInfo]chan<- struct{})
	set := file.NewChainSet(0, time.Second*10)
	set.Run()
	for i, api := range conf.APIs {
		app, _ := set.NewChain(api.Name, api.Domain+api.Route, false)
		// creating the jobs
		exec, ch, err := job.NewJob("monitor", app, &api)
		if err != nil {
			require.FailNow(t, "error creating # %d job: %s\n", i, err)
			continue
		}
		// launching the jobs
		go exec.Execute()
		jobs[exec.Info()] = ch
	}
	ctx, cancel := context.WithCancel(context.Background())
	scheduler := New(jobs)
	if len(scheduler.timeline) != len(jobs) {
		t.Errorf("error creating job : number of jobs in timeline is inaccurate")
	}
	go scheduler.Run(ctx)
	time.Sleep(time.Second * 30)
	cancel()
	time.Sleep(time.Second * 2)
	err = os.RemoveAll("storage")
	require.NoError(t, err)
}
