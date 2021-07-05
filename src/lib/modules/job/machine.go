package job

import (
	"fmt"
	"net/url"
	"os"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
	"github.com/bench-routes/bench-routes/src/lib/modules/evaluate"
	"github.com/bench-routes/bench-routes/tsdb/file"
)

type machineJob struct {
	JobInfo
	app   file.Appendable
	sigCh chan struct{}
	host  string
}

func newMachineJob(app file.Appendable, api *config.API) (*machineJob, chan<- struct{}, error) {
	sig := make(chan struct{})
	newjob := &machineJob{
		app:   app,
		sigCh: sig,
		JobInfo: JobInfo{
			Name:        api.Name,
			Every:       api.Every,
			lastExecute: time.Now().Add(api.Every * -1),
		},
	}
	url, err := url.Parse(api.Domain)
	if err != nil {
		return nil, nil, err
	}
	newjob.host = url.Host
	return newjob, sig, nil
}

// Execute execute the machineJob.
func (mn *machineJob) Execute() {
	for range mn.sigCh {
		mn.JobInfo.writeTime()
		ping, jitter, err := evaluate.Machine(mn.host)
		if err != nil {
			fmt.Fprintf(os.Stderr, "job: %s: error: %s", mn.JobInfo.Name, err.Error())
		}
		pingVal := fmt.Sprintf("%v|%v|%v", ping.Max.Microseconds(), ping.Mean.Microseconds(), ping.Mean.Microseconds())
		jitterVal := fmt.Sprintf("%v", jitter.Value.Microseconds())
		mn.app.Append(file.NewBlock("ping", pingVal))
		mn.app.Append(file.NewBlock("jitter", jitterVal))
	}
}

// Abort aborts the running job.
func (mn *machineJob) Abort() {
	close(mn.sigCh)
}

// Info returns the jobInfo of the job.
func (mn *machineJob) Info() *JobInfo {
	return &mn.JobInfo
}
