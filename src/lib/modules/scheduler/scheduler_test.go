package scheduler

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/bench-routes/bench-routes/src/lib/modules/job"
	file "github.com/bench-routes/bench-routes/tsdb"
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
		app, _ := set.NewChain(api.Name+"_monitor", api.Protocol+api.Domain+api.Route, false)
		// creating the jobs
		exec, ch, err := job.NewJob("monitor", app, nil, &api)
		if err != nil {
			require.FailNow(t, "error creating # %d job: %s\n", i, err)
			continue
		}
		// launching the jobs
		errCh := make(chan error)
		go func() {
			err, ok := <-errCh
			if !ok {
				return
			}
			if err != nil {
				fmt.Printf("job error : %v", err)
				return
			}
		}()
		go exec.Execute(errCh)
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
