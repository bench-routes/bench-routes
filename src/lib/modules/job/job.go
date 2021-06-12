package job

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

type Executable interface{
	Execute() 
	Abort() error
	Info() (error,jobInfo) 
}

type appendable interface{

}

type jobInfo struct {
    mux 			*sync.RWMutex
    name 			string
    every 			time.Duration
    lastExecute 	time.Time
}
 
type monitoringJob struct {
    jobInfo
    app 			appendable
    client 			*http.Client
    request 		*http.Request
}

type machineJob struct {
    jobInfo
    app 			appendable
    ping 			*ping.Pinger
}

func(mn monitoringJob) Execute(c <-chan struct{}) {
	for range c{
		//do monitoring

	}
}

func newMonitoringJob(app appendable, name string, every time.Duration) *monitoringJob{
	return &monitoringJob{
		app: app,
		jobInfo: jobInfo{
			name : name,
			every: every,
			lastExecute : time.Now().Add(every*-1),
		},
	}
}

func newMachineJob(app appendable, name, domainOrUrl string, every time.Duration) *machineJob{
	return &machineJob{
		app: app,
		jobInfo: jobInfo{
			name : name,
			every: every,
			lastExecute : time.Now().Add(every*-1),
		},
	}
}

func(mn monitoringJob) Abort(c chan<- struct{}) {
	//do monitoring
	close(c)
}

func(mn monitoringJob) Info() jobInfo{
	//do monitoring

	return jobInfo{};
}

func(mn machineJob) Execute(c <-chan struct{}){
	for range c{
		//do monitoring
		
	}
}

func(mn machineJob) Abort(c chan<- struct{}) {
	//do monitoring
	close(c)
}

func(mn machineJob) Info() jobInfo{
	//do monitoring
	return jobInfo{};
}