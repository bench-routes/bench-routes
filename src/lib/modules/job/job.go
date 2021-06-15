package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
	"github.com/bench-routes/bench-routes/tsdb/file"
)

type Executable interface{
	Execute() 
	Abort() 
	Info() *jobInfo
}

type jobInfo struct {
    // mux 			*sync.RWMutex
    name 			string
    every 			time.Duration
    lastExecute 	time.Time
}
 
type monitoringJob struct {
    jobInfo
    app 			*file.Appendable
	sigCh			chan struct{}
    client 			*http.Client
    // request 		*http.Request
	api				*config.API
}

type machineJob struct {
    jobInfo
    app 			*file.Appendable
	sigCh			chan struct{}
    // ping 			*ping.Pinger
}

func NewJob(typ string,c chan struct{},api *config.API) (Executable,error){
	var app *file.Appendable
	switch typ {
	case "machine":
		job,err := newMachineJob(app,c,api)
		if err != nil {
			return nil, fmt.Errorf("error creating job : %w",err)
		}
		return job,nil
	case "monitor":
		job,err := newMonitoringJob(app,c,api)
		if err != nil {
			return nil, fmt.Errorf("error creating job : %w",err)
		}
		return job,nil
	default:
		return nil,fmt.Errorf("`typ` provided is invalid")
	}	
}

func newMonitoringJob(app *file.Appendable,c chan struct{},api *config.API) (*monitoringJob,error){
	newJob := &monitoringJob{
		app: app,
		sigCh : c,
		jobInfo: jobInfo{
			name : api.Name,
			every: api.Every,
			lastExecute : time.Now().Add(api.Every*-1),
		},
		client: &http.Client{
			Timeout: time.Second*5,
		},
		api: api,
	}

	return newJob,nil
}

//Execute the monitoring function of the job when signal is passed to the chan sigCh
func(mn *monitoringJob) Execute() {
	for range mn.sigCh{
		log.Println("Inside Execute")
		request,err := newReq(mn.api)
		// request,err := http.NewRequest(strings.ToUpper(mn.api.Method),mn.api.Domain,nil)
		// request,err := http.NewRequest(strings.ToUpper("GET"),"https://www.google.com",nil)
		if err != nil {
			panic(err)
		} 
		stamp := time.Now();
 		res,err := mn.client.Do(request)
		if err != nil {
			fmt.Println(fmt.Errorf("error in sending request: %w",err));
		}
		resDelay := time.Since(stamp);
		resBody,err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(fmt.Errorf("error in reading response body: %w",err));
		}
		res.Body.Close()
		fmt.Printf("resDelay : %v\n",resDelay)
		fmt.Printf("resLength : %v\n",len(resBody))
	}
}

//Stop monitoring particular job
func(mn *monitoringJob) Abort(){
	close(mn.sigCh)
}

//jobInfo of the monitored job
func(mn *monitoringJob) Info() *jobInfo{
	return &mn.jobInfo;
}

//Creates http client and request for the job
func newReq(api *config.API) (*http.Request,error){
	body,err := json.Marshal(api.Body)
	if err != nil {
		return nil,fmt.Errorf("error marshalling : %w",err)
	}
	log.Println("Body : ",api.Body)
	//Creating http request for job
	// request,err := http.NewRequest(strings.ToUpper(api.Method),api.Domain,nil) 
	request,err := http.NewRequest(strings.ToUpper(api.Method),api.Domain+api.Route,bytes.NewBuffer(body)) 
	if err != nil {
		return nil,fmt.Errorf("error making http request : %w",err)
	}
	// adding headers to the request
	for k,v := range api.Headers{
		request.Header.Set(k,v);
	}
	return request,nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func newMachineJob(app *file.Appendable,c chan struct{},api *config.API) (*machineJob,error){
	return &machineJob{
		app: app,
		sigCh : c,
		jobInfo: jobInfo{
			name : api.Name,
			every: api.Every,
			lastExecute : time.Now().Add(api.Every*-1),
		},
	},nil
}


//Executing the
func(mn *machineJob) Execute(){
	for range mn.sigCh{
		//do monitoring
		
	}
}

func(mn *machineJob) Abort() {
	//do monitoring
	close(mn.sigCh)
}

func(mn *machineJob) Info() *jobInfo{
	//do monitoring
	return &mn.jobInfo;
}