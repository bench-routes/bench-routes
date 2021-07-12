package job

import (
	"fmt"
	"sync"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/bench-routes/bench-routes/tsdb/file"
)

// Executable is an interface that is implemented by machineJob and monitoringJob.
type Executable interface {
	Execute(chan<- error)
	Abort()
	Info() *JobInfo
}

// NewJob creates a new job based on the typ.
func NewJob(typ string, app file.Appendable, api *config.API) (Executable, chan<- struct{}, error) {
	switch typ {
	case "machine":
		job, ch, err := newMachineJob(app, api)
		if err != nil {
			return nil, nil, fmt.Errorf("creating machine job : %w", err)
		}
		return job, ch, nil
	case "monitor":
		job, ch, err := newMonitoringJob(app, api)
		if err != nil {
			return nil, nil, fmt.Errorf("creating monitoring job : %w", err)
		}
		return job, ch, nil
	default:
		return nil, nil, fmt.Errorf("`typ` provided is invalid")
	}
}

// JobInfo contains information about the job.
type JobInfo struct {
	mux         sync.RWMutex
	Name        string
	Every       time.Duration
	lastExecute time.Time
}

// ReadTime reads lastExecute of the JobInfo.
func (j *JobInfo) ReadTime() time.Time {
	j.mux.RLock()
	t := j.lastExecute
	j.mux.RUnlock()
	return t
}

// writeTime writes lastExecute of the JobInfo.
func (j *JobInfo) writeTime() {
	j.mux.Lock()
	j.lastExecute = time.Now()
	j.mux.Unlock()
}
