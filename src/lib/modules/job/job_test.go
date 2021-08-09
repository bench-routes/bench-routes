package job

import (
	"fmt"
	"os"
	"testing"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config"
	file "github.com/bench-routes/bench-routes/tsdb"
	"github.com/stretchr/testify/require"
)

var testapis []config.API = []config.API{
	{
		Name:     "API_1",
		Every:    time.Second * 2,
		Protocol: "http",
		Domain:   "www.youtube.com",
		Route:    "/watch",
		Method:   "GET",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Params: map[string]string{
			"v": "qqNenU-s",
		},
	},
	{
		Name:     "API_2",
		Every:    time.Second * 3,
		Protocol: "http",
		Domain:   "reqres.in",
		Route:    "/api/users",
		Method:   "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]string{
			"name": "tushar",
			"job":  "SDE",
		},
	},
}

func TestMonitorJob(t *testing.T) {
	set := file.NewChainSet(0, time.Second*10)
	set.Run()
	var paths []string

	// Running test for all the given apis.
	for index, api := range testapis {
		fmt.Printf("testing ResDelay and resLength for %s\n", api.Name)
		app, utils := set.NewChain(api.Name+"_monitor", api.Protocol+api.Domain+api.Route, true)
		paths = append(paths, utils.Path())
		exec, ch, err := NewJob("monitor", app, nil, &api)
		if err != nil {
			require.FailNow(t, "error creating %d # job : %s", index, err)
		}
		if exec == nil {
			require.FailNow(t, "error creating %d # job : returns nil", index)
		}
		if exec.Info().Every != api.Every || exec.Info().Name != api.Name {
			require.FailNow(t, "error creating %d # job : jobInfo is not correct", index)
		}
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
		for i := 0; i < 1; i++ {
			ch <- struct{}{}
			time.Sleep(exec.Info().Every)
		}
		exec.Abort()
	}

	time.Sleep(time.Second * 10)
	// Deleting all the test files generated.
	for _, path := range paths {
		err := os.Remove(path)
		require.NoError(t, err)
	}
	err := os.Remove("testdata")
	require.NoError(t, err)
}

func TestMachineJob(t *testing.T) {
	set := file.NewChainSet(0, time.Second*10)
	set.Run()
	var paths []string

	// Running tests for all the apis.
	for index, api := range testapis {
		fmt.Printf("testing ResDelay and resLength for %s\n", api.Name)
		appPing, utils1 := set.NewChain(api.Domain+"_ping", api.Protocol+"://"+api.Domain+api.Route, true)
		appJitter, utils2 := set.NewChain(api.Domain+"_jitter", api.Protocol+"://"+api.Domain+api.Route, true)
		paths = append(paths, utils1.Path())
		paths = append(paths, utils2.Path())
		exec, ch, err := NewJob("machine", appPing, appJitter, &api)
		if err != nil {
			require.FailNow(t, "error creating %d # job : %s", index, err)
		}
		if exec == nil {
			require.FailNow(t, "error creating %d # job : returns nil", index)
		}
		if exec.Info().Every != api.Every || exec.Info().Name != api.Name {
			require.FailNow(t, "error creating %d # job : jobInfo is not correct", index)
		}
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
		for i := 0; i < 1; i++ {
			ch <- struct{}{}
			time.Sleep(exec.Info().Every)
		}
		exec.Abort()
	}

	time.Sleep(time.Second * 10)
	// Deleting all the test files generated.
	for _, path := range paths {
		err := os.Remove(path)
		require.NoError(t, err)
	}
	err := os.Remove("testdata")
	require.NoError(t, err)
}
