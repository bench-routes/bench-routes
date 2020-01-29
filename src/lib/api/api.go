package api

import (
	"encoding/json"
	"net/http"

	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/parser"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zcalusic/sysinfo"
)

const (
	testFilesDir = "tests/"
)

// API type for implementing the API interface.
type API struct {
	SysInfo   sysinfo.SysInfo `json:"systemInformation"`
	RequestIP string          `json:"requestIPAddress"`
	Data      interface{}     `json:"data"`
}

// New returns the API type for implementing the API interface.
func New() *API {
	tmp := sysinfo.SysInfo{}
	tmp.GetSysInfo()
	return &API{
		SysInfo: tmp,
	}
}

// Home handles the requests for the home page.
func (a *API) Home(w http.ResponseWriter, r *http.Request) {
	msg := "ping from " + r.RemoteAddr + ", sent pong in response"
	logger.Terminal(msg, "p")
}

// TestTemplate handles template related request for testing purposes.
func (a *API) TestTemplate(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, testFilesDir+"bench-routes-socket-tester.html")
}

// ServiceState handles requests related to the state of the services in
// the application.
func (a *API) ServiceState(w http.ResponseWriter, r *http.Request) {
	p := parser.New(utils.ConfigurationFilePath)
	p = p.Refresh()

	a.Data = struct {
		Ping       string `json:"ping"`
		FloodPing  string `json:"floodping"`
		Jitter     string `json:"jitter"`
		Monitoring string `json:"monitoring"`
	}{
		Ping:       p.Config.UtilsConf.ServicesSignal.Ping,
		FloodPing:  p.Config.UtilsConf.ServicesSignal.FloodPing,
		Jitter:     p.Config.UtilsConf.ServicesSignal.Jitter,
		Monitoring: p.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring,
	}
	a.setRequestIPAddress(r)

	a.send(w, a.marshalled())
}

func (a *API) setRequestIPAddress(r *http.Request) {
	a.RequestIP = r.RemoteAddr
}

func (a *API) marshalled() []byte {
	js, err := json.Marshal(*a)
	if err != nil {
		panic(err)
	}

	return js
}

func (a *API) send(w http.ResponseWriter, data []byte) {
	if _, err := w.Write(data); err != nil {
		panic(err)
	}
}
