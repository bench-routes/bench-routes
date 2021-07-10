package module

import (
	"fmt"

	config "github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/bench-routes/bench-routes/tsdb/file"
)

// Runnable is an interface that is implmented by Machine and Monitor.
type Runnable interface {
	Run()
	Reload(*config.Config) error
	Stop()
}

type Type uint8

const (
	MachineType Type = iota
	MonitorType
)

// New returns a Runnable interface to implement machine and monitoring jobs.
func New(typ Type, chainSet *file.ChainSet, errCh chan<- error) (Runnable, error) {
	switch typ {
	case MachineType:
		module, err := newMachineModule(chainSet, errCh)
		if err != nil {
			return nil, fmt.Errorf("creating machine module : %w", err)
		}
		return module, nil
	case MonitorType:
		module, err := newMonitorModule(chainSet, errCh)
		if err != nil {
			return nil, fmt.Errorf("creating monitoring module : %w", err)
		}
		return module, nil
	default:
		return nil, fmt.Errorf("`typ` provided is invalid")
	}
}
