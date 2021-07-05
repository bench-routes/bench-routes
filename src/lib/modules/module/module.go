package module

import (
	"fmt"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
)

// Runnable is an interface that is implmented by Machine and Monitor.
type Runnable interface {
	Run()
	Reload(*config.Config, chan<- error)
	Stop()
}

// New returns a Runnable interface to implment machine and monitoring jobs.
func New(typ string) (Runnable, error) {
	switch typ {
	case "machine":
		module, err := newMachineModule()
		if err != nil {
			return nil, fmt.Errorf("creating machine module : %w", err)
		}
		return module, nil
	case "monitor":
		module, err := newMonitorModule()
		if err != nil {
			return nil, fmt.Errorf("creating monitoring module : %w", err)
		}
		return module, nil
	default:
		return nil, fmt.Errorf("`typ` provided is invalid")
	}
}
