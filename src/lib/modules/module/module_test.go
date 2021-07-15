package module

import (
	"fmt"
	"os"
	"testing"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/bench-routes/bench-routes/tsdb/file"
	"github.com/stretchr/testify/require"
)

func TestMachineModule(t *testing.T) {
	set := file.NewChainSet(0, time.Second*2)
	set.Run()
	errCh := make(chan error)
	go func() {
		err, ok := <-errCh
		if !ok {
			return
		}
		if err != nil {
			fmt.Printf("module error : %v", err)
			return
		}
	}()
	module, err := New(MachineType, set, errCh)
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	go module.Run()

	conf, err := config.New("./testdata/config.yml")
	if err != nil {
		t.Fatalf("error loading config: %s", err)
	}

	go module.Reload(conf)

	time.Sleep(time.Second * 30)
	module.Stop()
	time.Sleep(time.Second * 2)
	err = os.RemoveAll("storage")
	require.NoError(t, err)
}

func TestMonitorModule(t *testing.T) {
	set := file.NewChainSet(0, time.Second*2)
	set.Run()
	errCh := make(chan error)
	go func() {
		err, ok := <-errCh
		if !ok {
			return
		}
		if err != nil {
			fmt.Printf("module error : %v", err)
			return
		}
	}()
	module, err := New(MonitorType, set, errCh)
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	go module.Run()

	conf, err := config.New("./testdata/config.yml")
	if err != nil {
		t.Fatalf("error loading config: %s", err)
	}
	go module.Reload(conf)

	time.Sleep(time.Second * 30)
	module.Stop()
	time.Sleep(time.Second * 2)
	err = os.RemoveAll("storage")
	require.NoError(t, err)
}
