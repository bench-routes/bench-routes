package job

import (
	"fmt"
	"sync"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
	"github.com/bench-routes/bench-routes/tsdb/file"
)

// Executable is an interface that is implemented by machineJob and monitoringJob.
type Executable interface {
	Execute()
	Abort()
	Info() *JobInfo
}

// NewJob creates a new job based on the typ.
func NewJob(typ string, app file.Appendable, c chan struct{}, api *config.API) (Executable, error) {
	switch typ {
	case "machine":
		job, err := newMachineJob(app, c, api)
		if err != nil {
			return nil, fmt.Errorf("error creating job : %w", err)
		}
		return job, nil
	case "monitor":
		job, err := newMonitoringJob(app, c, api)
		if err != nil {
			return nil, fmt.Errorf("error creating job : %w", err)
		}
		return job, nil
	default:
		return nil, fmt.Errorf("`typ` provided is invalid")
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
