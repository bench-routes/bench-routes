package job

import (
	"fmt"
	"net/url"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/bench-routes/bench-routes/src/lib/modules/evaluate"
	file "github.com/bench-routes/bench-routes/tsdb"
)

type machineJob struct {
	JobInfo
	appPing   file.Appendable
	appJitter file.Appendable
	sigCh     chan struct{}
	urlPath   string
}

func newMachineJob(appPing file.Appendable, appJitter file.Appendable, api *config.API) (*machineJob, chan<- struct{}, error) {
	sig := make(chan struct{})
	newjob := &machineJob{
		appPing:   appPing,
		appJitter: appJitter,
		sigCh:     sig,
		JobInfo: JobInfo{
			Name:        api.Name,
			Every:       api.Every,
			lastExecute: time.Now().Add(api.Every * -1),
		},
	}
	urlParsed, err := url.Parse(api.Domain)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing api.Domain: %w", err)
	}
	newjob.urlPath = urlParsed.Path
	return newjob, sig, nil
}

// Execute execute the machineJob.
func (mn *machineJob) Execute(errCh chan<- error) {
	for range mn.sigCh {
		mn.JobInfo.writeTime()
		ping, jitter, err := evaluate.Machine(mn.urlPath)
		if err != nil {
			errCh <- fmt.Errorf("job: %s: %w", mn.JobInfo.Name, err)
			return
		}
		pingVal := fmt.Sprintf("%v|%v|%v", ping.Max.Microseconds(), ping.Mean.Microseconds(), ping.Min.Microseconds())
		jitterVal := fmt.Sprintf("%v", jitter.Value.Microseconds())
		mn.appPing.Append(file.NewBlock("ping", pingVal))
		mn.appJitter.Append(file.NewBlock("jitter", jitterVal))
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
