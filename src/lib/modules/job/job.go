package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
	"github.com/bench-routes/bench-routes/src/lib/modules/evaluate"
	"github.com/bench-routes/bench-routes/tsdb/file"
)

// Executable is an interface that is implemented by machineJob and monitoringJob
type Executable interface {
	Execute()
	Abort()
	Info() *JobInfo
}

type JobInfo struct {
	mux         sync.RWMutex
	Name        string
	Every       time.Duration
	lastExecute time.Time
}

// NewJob creates a new job based on the typ
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

type monitoringJob struct {
	JobInfo
	app     file.Appendable
	sigCh   chan struct{}
	client  *http.Client
	request *http.Request
}

func newMonitoringJob(app file.Appendable, c chan struct{}, api *config.API) (*monitoringJob, error) {
	newJob := &monitoringJob{
		app:   app,
		sigCh: c,
		JobInfo: JobInfo{
			Name:        api.Name,
			Every:       api.Every,
			lastExecute: time.Now().Add(api.Every * -1),
		},
		client: &http.Client{
			Timeout: time.Second * 5,
		},
		request: &http.Request{},
	}
	if err := newReq(newJob, api); err != nil {
		return nil, fmt.Errorf("error making http request : %w", err)
	}
	return newJob, nil
}

// Execute executes the monitoringJob
func (mn *monitoringJob) Execute() {
	for range mn.sigCh {
		mn.JobInfo.writeTime()
		response, err := evaluate.Monitor(mn.client, mn.request)
		if err != nil {
			fmt.Println(fmt.Errorf("%s", err))
		}
		val := fmt.Sprintf("%v|%v", response.Delay.Microseconds(), strconv.Itoa(response.Length))
		fmt.Println(val)
		mn.app.Append(file.NewBlock("monitoring", val))
	}
}

// Aborts aborts the running job
func (mn *monitoringJob) Abort() {
	close(mn.sigCh)
}

// Info returns jobInfo of the job
func (mn *monitoringJob) Info() *JobInfo {
	return &mn.JobInfo
}

// newReq creates http client and request for the job
func newReq(job *monitoringJob, api *config.API) error {
	var body []byte
	var err error
	if len(api.Body) != 0 {
		body, err = json.Marshal(api.Body)
		if err != nil {
			return fmt.Errorf("error marshalling : %w", err)
		}
	}
	// creating http request for job
	job.request, err = http.NewRequest(strings.ToUpper(api.Method), api.Domain+api.Route, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error making http request : %w", err)
	}
	// adding headers to the request
	for k, v := range api.Headers {
		job.request.Header.Set(k, v)
	}
	return nil
}

type machineJob struct {
	JobInfo
	app   file.Appendable
	sigCh chan struct{}
	host  string
}

func newMachineJob(app file.Appendable, c chan struct{}, api *config.API) (*machineJob, error) {
	newjob := &machineJob{
		app:   app,
		sigCh: c,
		JobInfo: JobInfo{
			Name:        api.Name,
			Every:       api.Every,
			lastExecute: time.Now().Add(api.Every * -1),
		},
	}
	url, err := url.Parse(api.Domain)
	if err != nil {
		return nil, err
	}
	newjob.host = url.Host
	return newjob, nil
}

// Execute execute the machineJob
func (mn *machineJob) Execute() {
	for range mn.sigCh {
		mn.JobInfo.writeTime()
		ping, jitter, err := evaluate.Machine(mn.host)
		if err != nil {
			fmt.Println(fmt.Errorf("%s", err))
		}
		pingval := fmt.Sprintf("%v|%v|%v", ping.Max.Microseconds(), ping.Mean.Microseconds(), ping.Mean.Microseconds())
		jitterval := fmt.Sprintf("%v", jitter.Value.Microseconds())
		fmt.Println(pingval)
		fmt.Println(jitterval)
		mn.app.Append(file.NewBlock("ping", pingval))
		mn.app.Append(file.NewBlock("jitter", jitterval))
	}
}

// Abort aborts the running job
func (mn *machineJob) Abort() {
	close(mn.sigCh)
}

// Info returns the jobInfo of the job
func (mn *machineJob) Info() *JobInfo {
	return &mn.JobInfo
}

// ReadTime reads lastExecute of the JobInfo
func (j *JobInfo) ReadTime() time.Time {
	j.mux.RLock()
	t := j.lastExecute
	j.mux.RUnlock()
	return t
}

// writeTime writes lastExecute of the JobInfo
func (j *JobInfo) writeTime() {
	j.mux.Lock()
	j.lastExecute = time.Now()
	j.mux.Unlock()
}
