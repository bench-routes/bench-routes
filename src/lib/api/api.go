package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/parser"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/src/lib/utils/brt"
	"github.com/zairza-cetb/bench-routes/tsdb"
	"github.com/zairza-cetb/bench-routes/tsdb/querier"
)

const (
	testFilesDir = "tests/"
	uiPathV1     = "ui-builds/v1.0/index.html"
)

// API type for implementing the API interface.
type API struct {
	RequestIP string
	Data      interface{}
	Matrix    *utils.BRmap
}

// New returns the API type for implementing the API interface.
func New(matrix *utils.BRmap) *API {
	return &API{
		Matrix: matrix,
	}
}

// Home handles the requests for the home page.
func (a *API) Home(w http.ResponseWriter, r *http.Request) {
	msg := "ping from " + r.RemoteAddr + ", sent pong in monitor"
	logger.Terminal(msg, "p")
}

// UIv1 serves the v1.0 version of user-interface of bench-routes.
// ui-builds/v1.0 is served through this.
func (a *API) UIv1(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir(uiPathV1))
}

// TestTemplate handles template related request for testing purposes.
func (a *API) TestTemplate(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, testFilesDir+"bench-routes-socket-tester.html")
}

// ServiceState handles requests related to the state of the Services in
// the application.
func (a *API) ServiceState(w http.ResponseWriter, r *http.Request) {
	p := parser.New(utils.ConfigurationFilePath)
	p.Refresh()

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

// RoutesSummary handles requests related to summarized-configuration details.
func (a *API) RoutesSummary(w http.ResponseWriter, r *http.Request) {
	p := parser.New(utils.ConfigurationFilePath)
	p.Refresh()

	var servicesRoutes, monitoringRoutes []string
	for _, r := range p.Config.Routes {
		servicesRoutes = append(servicesRoutes, r.URL)
		monitoringRoutes = append(monitoringRoutes, r.Method+": "+r.URL+"/"+r.Route)
	}

	a.Data = struct {
		TestServicesRoutes []string `json:"testServicesRoutes"`
		MonitoringRoutes   []string `json:"monitoringRoutes"`
	}{
		TestServicesRoutes: servicesRoutes,
		MonitoringRoutes:   monitoringRoutes,
	}
	a.setRequestIPAddress(r)
	a.send(w, a.marshalled())
}

// Query forms the query handler for querying over the time-series data.
func (a *API) Query(w http.ResponseWriter, r *http.Request) {
	var (
		startTimestamp, endTimestamp int64
		err                          error
	)
	timeSeriesPath := r.URL.Query().Get("timeSeriesPath")

	startTimestampStr := r.URL.Query().Get("startTimestamp")
	if startTimestampStr == "" {
		startTimestamp = int64(math.MaxInt64)
	} else {
		startTimestamp, err = strconv.ParseInt(startTimestampStr, 10, 64)
		if err != nil {
			logger.Terminal(fmt.Errorf("error in startTimestamp: %s", err.Error()).Error(), "p")
		}
	}

	endTimestampStr := r.URL.Query().Get("endTimestamp")
	if endTimestampStr == "" {
		endTimestamp = int64(math.MinInt64)
	} else {
		endTimestamp, err = strconv.ParseInt(endTimestampStr, 10, 64)
		if err != nil {
			logger.Terminal(fmt.Errorf("error in endTimestamp: %s", err.Error()).Error(), "p")
		}
	}

	// condition: only for bench-routes as per the design
	//
	// path should be in syntax: <DBPath>/<ofType>/chunk_<middle>_<url>.json -> non-system metric
	// %s/system.json -> system metric

	// verify if chain path exists
	timeSeriesPath = timeSeriesPath + tsdb.TSDBFileExtension
	if ok := tsdb.VerifyChainPathExists(timeSeriesPath); !ok {
		a.send(w, []byte("INVALID_PATH"))
		return
	}

	// TODO: we do not capture the block streams in memory while querying yet. They are captured only when flushed.
	// TODO: consider cmap while querying for fetching latest blocks after shifting tsdb to binary.

	qry := querier.New(timeSeriesPath, "")
	query := qry.QueryBuilder()
	query.SetRange(startTimestamp, endTimestamp)
	a.Data = query.ExecWithoutEncode()
	a.send(w, a.marshalled())
}

// SendMatrix responds by sending the multi-dimensional data (called matrix)
// dependent on a route name as in matrix key.
func (a *API) SendMatrix(w http.ResponseWriter, r *http.Request) {
	routeNameMatrix := r.URL.Query().Get("routeNameMatrix")
	if _, ok := (*a.Matrix)[routeNameMatrix]; !ok {
		a.Data = "ROUTE_NAME_NOT_IN_MATRIX"
		a.send(w, a.marshalled())
		return
	}
	curr := time.Now().UnixNano()
	from := curr - (brt.Hour * 6)

	matrix := (*a.Matrix)[routeNameMatrix]
	parallelQueryExec := func(path string, c chan querier.QueryResponse) {
		qry := querier.New(path, "")
		query := qry.QueryBuilder()
		query.SetRange(curr, from)
		c <- query.ExecWithoutEncode()
	}
	chans := []chan querier.QueryResponse{
		make(chan querier.QueryResponse),
		make(chan querier.QueryResponse),
		make(chan querier.QueryResponse),
		make(chan querier.QueryResponse),
	}
	go parallelQueryExec(matrix.PingChain.Path, chans[0])
	go parallelQueryExec(matrix.FPingChain.Path, chans[1])
	go parallelQueryExec(matrix.JitterChain.Path, chans[2])
	go parallelQueryExec(matrix.MonitorChain.Path, chans[3])
	matrixResponse := map[string]querier.QueryResponse{
		"ping":    <-chans[0],
		"fping":   <-chans[1],
		"jitter":  <-chans[2],
		"monitor": <-chans[3],
	}
	a.Data = matrixResponse
	a.send(w, a.marshalled())
}

// TSDBPathDetails responds with the path details that will be used for
// passing into the querier's timeSeriesPath.
func (a *API) TSDBPathDetails(w http.ResponseWriter, _ *http.Request) {
	var chainDetails []utils.ResponseTSDBChains
	for _, v := range *a.Matrix {
		chainDetails = append(chainDetails, utils.ResponseTSDBChains{
			Name: v.Domain,
			Path: utils.ChainPath{
				Ping:    v.PingChain.Path,
				Jitter:  v.JitterChain.Path,
				Fping:   v.FPingChain.Path,
				Monitor: v.MonitorChain.Path,
			},
		})
	}
	a.send(w, a.marshal(chainDetails))
}

func (a *API) setRequestIPAddress(r *http.Request) {
	a.RequestIP = r.RemoteAddr
}

func (a *API) marshalled() []byte {
	response := struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}{
		Status: a.RequestIP,
		Data:   a.Data,
	}
	js, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return js
}

func (a *API) marshal(data interface{}) []byte {
	js, err := json.Marshal(data)
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
