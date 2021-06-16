package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
	"github.com/bench-routes/bench-routes/tsdb/file"
	"github.com/go-ping/ping"
)

type Executable interface {
	Execute()
	Abort()
	Info() *jobInfo
}

type jobInfo struct {
	mux         sync.RWMutex
	name        string
	every       time.Duration
	lastExecute time.Time
}

type monitoringJob struct {
	jobInfo
	app     *file.Appendable
	sigCh   chan struct{}
	client  *http.Client
	request *http.Request
}

type machineJob struct {
	jobInfo
	app   *file.Appendable
	sigCh chan struct{}
	host  string
}

func NewJob(typ string, c chan struct{}, api *config.API) (Executable, error) {
	var app *file.Appendable
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

func newMonitoringJob(app *file.Appendable, c chan struct{}, api *config.API) (*monitoringJob, error) {
	newJob := &monitoringJob{
		app:   app,
		sigCh: c,
		jobInfo: jobInfo{
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

//Execute the monitoring function of the job when signal is passed to the chan sigCh
func (mn *monitoringJob) Execute() {
	for range mn.sigCh {
		log.Println("Inside Execute")
		stamp := time.Now()
		res, err := mn.client.Do(mn.request)
		if err != nil {
			fmt.Println(fmt.Errorf("error in sending request: %w", err))
		}
		resDelay := time.Since(stamp)

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(fmt.Errorf("error in reading response body: %w", err))
		}
		res.Body.Close()
		mn.jobInfo.mux.Lock()
		mn.jobInfo.lastExecute = time.Now()
		mn.jobInfo.mux.Unlock()
		fmt.Printf("resDelay : %v\n", resDelay)
		fmt.Printf("resLength : %v\n", len(resBody))
	}
}

//Stop monitoring particular job
func (mn *monitoringJob) Abort() {
	close(mn.sigCh)
}

// jobInfo of the monitored job
func (mn *monitoringJob) Info() *jobInfo {
	return &mn.jobInfo
}

// Creates http client and request for the job
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

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func newMachineJob(app *file.Appendable, c chan struct{}, api *config.API) (*machineJob, error) {
	newjob := &machineJob{
		app:   app,
		sigCh: c,
		jobInfo: jobInfo{
			name:        api.Name,
			every:       api.Every,
			lastExecute: time.Now().Add(api.Every * -1),
		},
		host: api.Domain,
	}
	return newjob, nil
}

//Executing the
func (mn *machineJob) Execute() {
	for range mn.sigCh {
		pinger, err := ping.NewPinger(mn.host)
		if err != nil {
			fmt.Println(fmt.Errorf("error creating ping : %w", err))
		}
		pinger.Count = 5
		var lastTime time.Duration
		var sum time.Duration
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
		}
		mn.jobInfo.mux.Lock()
		mn.jobInfo.lastExecute = time.Now()
		mn.jobInfo.mux.Unlock()
		stats := pinger.Statistics()
		fmt.Println("Ping : ", stats.AvgRtt)
		fmt.Println("Jitter : ", sum/time.Duration(pinger.Count-1))
	}
}

func (mn *machineJob) Abort() {
	//do monitoring
	close(mn.sigCh)
}

func (mn *machineJob) Info() *jobInfo {
	//do monitoring
	return &mn.jobInfo
}

func absDiff(a, b time.Duration) time.Duration {
	if a >= b {
		return a - b
	}
	return b - a
}
