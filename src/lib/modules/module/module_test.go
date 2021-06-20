package module

import (
	"testing"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
)

func TestModule(t *testing.T) {
	module, err := NewModule("monitor")
	if err != nil {
		t.Fatalf("error creating module %s", err)
	}

	go module.Run()

	conf, err := config.New("./testdata/config.yml")
	if err != nil {
		t.Fatalf("error loading config: %s", err)
	}

	module.Reload(conf)
	time.Sleep(time.Minute)
	module.Stop()
	time.Sleep(time.Second * 5)
}
