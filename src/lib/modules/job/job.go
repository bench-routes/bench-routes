package job

import (
	"fmt"
	"sync"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config"
	file "github.com/bench-routes/bench-routes/tsdb"
)

// Executable is an interface that is implemented by machineJob and monitoringJob.
type Executable interface {
	Execute(chan<- error)
	Abort()
	Info() *JobInfo
}

// NewJob creates a new job based on the typ.
func NewJob(typ string, app_1 file.Appendable, app_2 file.Appendable, api *config.API) (Executable, chan<- struct{}, error) {
	// Here we use app_1 and app_2 to represent two appenders. In case of
	// machine : app_1 will represent ping appendable and app_2 will represent jitter appendable
	// monitor : app_1 will represent monitor appendable and app_2 will be nil.
	switch typ {
	case "machine":
		job, ch, err := newMachineJob(app_1, app_2, api)
		if err != nil {
			return nil, nil, fmt.Errorf("creating machine job : %w", err)
		}
		return job, ch, nil
	case "monitor":
		job, ch, err := newMonitoringJob(app_1, api)
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
