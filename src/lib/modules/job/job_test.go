package job

import (
	"log"
	"testing"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
)

func TestJob(t *testing.T) {
	ch := make(chan struct{})
	api := config.API{
		Name:   "API_1",
		Every:  time.Second * 5,
		Domain: "https://reqres.in/api/users",
		// Route: "/watch?v=aaDaIGHTGT8",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]string{
			"name": "tushar",
			"job":  "SDE",
		},
	}
	exec, err := NewJob("monitor", ch, &api)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	go exec.Execute()
	ch <- struct{}{}
	log.Println(exec.Info().lastExecute)
	time.Sleep(time.Second * 5)

	ch <- struct{}{}
	log.Println(exec.Info().lastExecute)
	time.Sleep(time.Second * 5)

	ch <- struct{}{}
	log.Println(exec.Info().lastExecute)
	time.Sleep(time.Second * 5)

	ch <- struct{}{}
	log.Println(exec.Info().lastExecute)
	time.Sleep(time.Second * 5)

	exec.Abort()

}

func TestMachineJob(t *testing.T) {
	ch := make(chan struct{})
	api := config.API{
		Name:   "API_1",
		Every:  time.Second * 5,
		Domain: "reqres.in",
		Route: "/api/users",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]string{
			"name": "tushar",
			"job":  "SDE",
		},
	}
	exec, err := NewJob("machine", ch, &api)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	go exec.Execute()
	ch <- struct{}{}
	log.Println(exec.Info().lastExecute)
	time.Sleep(time.Second * 10)

	ch <- struct{}{}
	log.Println(exec.Info().lastExecute)
	time.Sleep(time.Second * 10)

	ch <- struct{}{}
	log.Println(exec.Info().lastExecute)
	time.Sleep(time.Second * 10)

	ch <- struct{}{}
	log.Println(exec.Info().lastExecute)
	time.Sleep(time.Second * 10)

	exec.Abort()

}
