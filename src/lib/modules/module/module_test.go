package module

import (
	"os"
	"testing"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/stretchr/testify/require"
)

func TestMachineModule(t *testing.T) {
	module, err := New("machine")
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	go module.Run()

	conf, err := config.New("./testdata/config.yml")
	if err != nil {
		t.Fatalf("error loading config: %s", err)
	}
	errCh := make(chan error)

	go module.Reload(conf, errCh)
	if err := <-errCh; err != nil {
		t.Fatalf("reloading module : %s", err)
	}
	time.Sleep(time.Second * 30)
	module.Stop()
	time.Sleep(time.Second * 2)
	err = os.RemoveAll("storage")
	require.NoError(t, err)
}

func TestMonitorModule(t *testing.T) {
	module, err := New("monitor")
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	go module.Run()

	conf, err := config.New("./testdata/config.yml")
	if err != nil {
		t.Fatalf("error loading config: %s", err)
	}
	errCh := make(chan error)
	go module.Reload(conf, errCh)
	if err := <-errCh; err != nil {
		t.Fatalf("reloading module : %s", err)
	}
	time.Sleep(time.Second * 40)
	module.Stop()
	time.Sleep(time.Second * 2)
	err = os.RemoveAll("storage")
	require.NoError(t, err)
}
