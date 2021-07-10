package job

import (
	"fmt"
	"os"
	"testing"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/bench-routes/bench-routes/tsdb/file"
	"github.com/stretchr/testify/require"
)

var testapis []config.API = []config.API{
	{
		Name:   "API_1",
		Every:  time.Second * 5,
		Domain: "https://www.youtube.com",
		Route:  "/watch",
		Method: "GET",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Params: map[string]string{
			"v": "qqNenU-s",
		},
	},
	{
		Name:   "API_2",
		Every:  time.Second * 10,
		Domain: "https://reqres.in",
		Route:  "/api/users",
		Method: "POST",
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
		app, utils := set.NewChain(api.Name, api.Domain+api.Route, true)
		paths = append(paths, utils.Path())
		exec, ch, err := NewJob("monitor", app, &api)
		if err != nil {
			require.FailNow(t, "error creating %d # job : %s", index, err)
		}
		if exec == nil {
			require.FailNow(t, "error creating %d # job : returns nil", index)
		}
		if exec.Info().Every != api.Every || exec.Info().Name != api.Name {
			require.FailNow(t, "error creating %d # job : jobInfo is not correct", index)
		}
		go exec.Execute()
		for i := 0; i < 3; i++ {
			ch <- struct{}{}
			time.Sleep(exec.Info().Every)
		}
		exec.Abort()
	}

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
		app, utils := set.NewChain(api.Name, api.Domain+api.Route, true)
		paths = append(paths, utils.Path())
		exec, ch, err := NewJob("machine", app, &api)
		if err != nil {
			require.FailNow(t, "error creating %d # job : %s", index, err)
		}
		if exec == nil {
			require.FailNow(t, "error creating %d # job : returns nil", index)
		}
		if exec.Info().Every != api.Every || exec.Info().Name != api.Name {
			require.FailNow(t, "error creating %d # job : jobInfo is not correct", index)
		}
		go exec.Execute()
		for i := 0; i < 3; i++ {
			ch <- struct{}{}
			time.Sleep(exec.Info().Every)
		}
		exec.Abort()
	}

	// Deleting all the test files generated.
	for _, path := range paths {
		err := os.Remove(path)
		require.NoError(t, err)
	}
	err := os.Remove("testdata")
	require.NoError(t, err)
}
