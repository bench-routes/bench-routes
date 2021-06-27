package job

import (
	"fmt"
	"testing"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
	"github.com/bench-routes/bench-routes/tsdb/file"
)

var testapis []config.API = []config.API{
	{
		Name:   "API_1",
		Every:  time.Second * 5,
		Domain: "https://www.youtube.com",
		Route:  "/watch?v=qqNenU-sQZI",
		Method: "GET",
		Headers: map[string]string{
			"Content-Type": "application/json",
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
	for index, api := range testapis {
		fmt.Printf("testing ResDelay and resLength for %s\n", api.Name)
		ch := make(chan struct{})
		app, _ := set.NewChain(api.Name, api.Domain+api.Route, true)
		exec, err := NewJob("monitor", app, ch, &api)
		if err != nil {
			t.Errorf("error creating %d # job : %s", index, err)
		}
		if exec == nil {
			t.Errorf("error creating %d # job : returns nil", index)
		}
		if exec.Info().Every != api.Every || exec.Info().Name != api.Name {
			t.Errorf("error creating %d # job : jobInfo is not correct", index)
		}
		go exec.Execute()
		for i := 0; i < 3; i++ {
			ch <- struct{}{}
			time.Sleep(time.Second * 5)
		}
		exec.Abort()
	}
}

func TestMachineJob(t *testing.T) {
	set := file.NewChainSet(0, time.Second*10)
	set.Run()
	for index, api := range testapis {
		fmt.Printf("testing ResDelay and resLength for %s\n", api.Name)
		ch := make(chan struct{})
		app, _ := set.NewChain(api.Name, api.Domain+api.Route, true)
		exec, err := NewJob("machine", app, ch, &api)
		if err != nil {
			t.Errorf("error creating %d # job : %s", index, err)
		}
		if exec == nil {
			t.Errorf("error creating %d # job : returns nil", index)
		}
		if exec.Info().Every != api.Every || exec.Info().Name != api.Name {
			t.Errorf("error creating %d # job : jobInfo is not correct", index)
		}
		go exec.Execute()
		for i := 0; i < 3; i++ {
			ch <- struct{}{}
			time.Sleep(time.Second * 10)
		}
		exec.Abort()
	}
}
