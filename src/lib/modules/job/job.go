package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
	"github.com/bench-routes/bench-routes/src/lib/filters"
	"github.com/bench-routes/bench-routes/tsdb/file"
	"github.com/go-ping/ping"
)

// Executable is an interface that is implemented by machineJob and monitoringJob
type Executable interface {
	Execute()
	Abort()
	Info() *JobInfo
}

type JobInfo struct {
	mux         sync.RWMutex
	name        string
	every       time.Duration
	lastExecute time.Time
}

type monitoringJob struct {
	JobInfo
	app     file.Appendable
	sigCh   chan struct{}
	client  *http.Client
	request *http.Request
}

type machineJob struct {
	JobInfo
	app   file.Appendable
	sigCh chan struct{}
	host  string
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

func newMonitoringJob(app file.Appendable, c chan struct{}, api *config.API) (*monitoringJob, error) {
	newJob := &monitoringJob{
		app:   app,
		sigCh: c,
		JobInfo: JobInfo{
			name:        api.Name,
			every:       api.Every,
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
		stamp := time.Now()
		res, err := mn.client.Do(mn.request)
		if err != nil {
			fmt.Println(fmt.Errorf("error in sending request: %w", err))
			continue
		}
		resDelay := time.Since(stamp)

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(fmt.Errorf("error in reading response body: %w", err))
			continue
		}
		res.Body.Close()
		mn.JobInfo.mux.Lock()
		mn.JobInfo.lastExecute = time.Now()
		mn.JobInfo.mux.Unlock()
		val := fmt.Sprintf("%s|%s", fmt.Sprint(resDelay), fmt.Sprint(len(resBody)))
		fmt.Println(val)
		
		mn.app.Append(file.NewBlock("job-monitoring",val))
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

func newMachineJob(app file.Appendable, c chan struct{}, api *config.API) (*machineJob, error) {
	newjob := &machineJob{
		app:   app,
		sigCh: c,
		JobInfo: JobInfo{
			name:        api.Name,
			every:       api.Every,
			lastExecute: time.Now().Add(api.Every * -1),
		},
		host: *filters.HTTPPingFilter(&api.Domain),
	}
	return newjob, nil
}

// Execute execute the machineJob
func (mn *machineJob) Execute() {
	for range mn.sigCh {
		pinger, err := ping.NewPinger(mn.host)
		if err != nil {
			fmt.Println(fmt.Errorf("error creating ping : %w", err))
			continue
		}
		pinger.Count = 5
		var lastTime time.Duration
		var sum time.Duration
		// Calculating jitter using ping values
		pinger.OnRecv = func(pkt *ping.Packet) {
			if lastTime != time.Second*0 {
				sum += absDiff(lastTime, pkt.Rtt)
				fmt.Println("lastTime : ", lastTime, " Current : ", pkt.Rtt, " Diff : ", absDiff(lastTime, pkt.Rtt))
			}
			lastTime = pkt.Rtt
		}
		// Runing the pinger
		if err := pinger.Run(); err != nil {
			fmt.Println(fmt.Errorf("error running ping : %w", err))
			continue
		}
		mn.JobInfo.mux.Lock()
		mn.JobInfo.lastExecute = time.Now()
		mn.JobInfo.mux.Unlock()
		stats := pinger.Statistics()
		fmt.Println("Ping : ", stats.AvgRtt)
		fmt.Println("Jitter : ", sum/time.Duration(pinger.Count-1))
	}
}

// Abort aborts the running job
func (mn *machineJob) Abort() {
	//do monitoring
	close(mn.sigCh)
}

// Info returns the jobInfo of the job
func (mn *machineJob) Info() *JobInfo {
	return &mn.JobInfo
}

func absDiff(a, b time.Duration) time.Duration {
	if a >= b {
		return a - b
	}
	return b - a
}
