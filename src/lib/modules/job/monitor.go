package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
	"github.com/bench-routes/bench-routes/src/lib/modules/evaluate"
	"github.com/bench-routes/bench-routes/tsdb/file"
)

type monitoringJob struct {
	JobInfo
	app     file.Appendable
	sigCh   chan struct{}
	client  *http.Client
	request *http.Request
}

// newMonitoringJob creates a new monitoringJob.
func newMonitoringJob(app file.Appendable, api *config.API) (*monitoringJob, chan<- struct{}, error) {
	sig := make(chan struct{})
	newJob := &monitoringJob{
		app:   app,
		sigCh: sig,
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
		return nil, nil, fmt.Errorf("error making http request : %w", err)
	}
	return newJob, sig, nil
}

// newReq creates http client and request for the job.
func newReq(job *monitoringJob, api *config.API) error {
	var (
		body []byte
		err  error
	)
	if len(api.Body) != 0 {
		body, err = json.Marshal(api.Body)
		if err != nil {
			return fmt.Errorf("error marshalling : %w", err)
		}
	}
	// creating http request for job.
	job.request, err = http.NewRequest(strings.ToUpper(api.Method), api.Domain+api.Route, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error making http request : %w", err)
	}
	// adding headers to the request.
	for k, v := range api.Headers {
		job.request.Header.Set(k, v)
	}
	// adding params to the request.
	params := url.Values{}
	for k, v := range api.Params {
		params.Add(k, v)
	}
	job.request.URL.RawQuery = params.Encode()
	return nil
}

// Execute executes the monitoringJob.
func (mn *monitoringJob) Execute() {
	for range mn.sigCh {
		mn.JobInfo.writeTime()
		response, err := evaluate.Monitor(mn.client, mn.request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "job: %s: error: %s", mn.JobInfo.Name, err.Error())
		}
		val := fmt.Sprintf("%v|%v|%v", response.Delay.Microseconds(), strconv.Itoa(response.Length), strconv.Itoa(response.Size))
		mn.app.Append(file.NewBlock("monitoring", val))
	}
}

// Aborts aborts the running job.
func (mn *monitoringJob) Abort() {
	close(mn.sigCh)
}

// Info returns jobInfo of the job.
func (mn *monitoringJob) Info() *JobInfo {
	return &mn.JobInfo
}
