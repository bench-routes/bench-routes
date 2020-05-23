package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/pprof"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/jitter"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/monitor"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/ping"
	"github.com/zairza-cetb/bench-routes/src/lib/parser"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/src/lib/utils/brt"
	"github.com/zairza-cetb/bench-routes/tsdb"
	"github.com/zairza-cetb/bench-routes/tsdb/querier"
)

const (
	testFilesDir = "tests/"
	uiPathV1     = "ui-builds/v1.0/index.html"
	uiPathV11    = "ui-builds/v1.1/"
)

// API type for implementing the API interface.
type API struct {
	ResponseStatus string
	Data, Services interface{}
	Matrix         *utils.BRmap
	config         *parser.YAMLBenchRoutesType
}

// New returns the API type for implementing the API interface.
func New(matrix *utils.BRmap, config *parser.YAMLBenchRoutesType, services interface{}) *API {
	return &API{
		Matrix:   matrix,
		Services: services,
		config:   config,
	}
}

// Register registers the routes with the mux router.
func (a *API) Register(router *mux.Router) {
	// static servings.
	{
		router.Handle("/", http.FileServer(http.Dir(uiPathV11)))
		router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(uiPathV11+"assets/"))))
		router.PathPrefix("/manifest.json").Handler(http.StripPrefix("/manifest.json", http.FileServer(http.Dir(uiPathV11+"/manifest.json"))))
		router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(uiPathV11+"static/"))))
	}
	// Pprof profiling routes.
	{
		/** Index responds with the pprof-formatted profile named by the request.
		For example, "/debug/pprof/heap" serves the "heap" profile.
		Index responds to a request for "/debug/pprof/" with an HTML page listing the available profiles. **/
		router.HandleFunc("/debug/pprof/", pprof.Index)
		/** Respective handlers for pprof.Index **/
		router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		router.Handle("/debug/pprof/block", pprof.Handler("block"))
		/** Cmdline responds with the running program's command line, with arguments separated by NUL bytes.
		The package initialization registers it as /debug/pprof/cmdline. **/
		router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		/** Profile responds with the pprof-formatted cpu profile.
		Profiling lasts for duration specified in seconds GET parameter,
		or for 30 seconds if not specified. The package initialization registers it as /debug/pprof/profile. **/
		router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		/** Symbol looks up the program counters listed in the request, responding
		with a table mapping program counters to function names.
		The package initialization registers it as /debug/pprof/symbol. **/
		router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	}
	router.HandleFunc("/br-live-check", a.Home)
	router.HandleFunc("/get-monitoring-services-state", a.GetMonitoringState)
	router.HandleFunc("/get-route-time-series", a.TSDBPathDetails)
	router.HandleFunc("/query-matrix", a.SendMatrix)
	router.HandleFunc("/query", a.Query)
	router.HandleFunc("/routes-summary", a.RoutesSummary)
	router.HandleFunc("/service-state", a.ServiceState)
	router.HandleFunc("/test", a.TestTemplate)
	router.HandleFunc("/update-monitoring-services-state", a.UpdateMonitoringServicesState)
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

	qry := querier.New(timeSeriesPath, "", querier.TypeRange)
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
	var matrixResponse map[string]querier.QueryResponse
	chans := []chan querier.QueryResponse{
		make(chan querier.QueryResponse),
		make(chan querier.QueryResponse),
		make(chan querier.QueryResponse),
	}
	startTimestampStr := r.URL.Query().Get("startTimestamp")
	endTimestampStr := r.URL.Query().Get("endTimestamp")
	parallelQueryExec := func(path string, Type uint8, curr, from int64, c chan querier.QueryResponse) {
		qry := querier.New(path, "", Type)
		query := qry.QueryBuilder()
		query.SetRange(curr, from)
		c <- query.ExecWithoutEncode()
	}
	matrix := (*a.Matrix)[routeNameMatrix]
	if startTimestampStr == "" && endTimestampStr == "" {
		curr := time.Now().UnixNano()
		from := curr - (brt.Minute * 20)

		go parallelQueryExec(matrix.PingChain.Path, querier.TypeFirst, curr, from, chans[0])
		go parallelQueryExec(matrix.JitterChain.Path, querier.TypeFirst, curr, from, chans[1])
		go parallelQueryExec(matrix.MonitorChain.Path, querier.TypeFirst, curr, from, chans[2])
		matrixResponse = map[string]querier.QueryResponse{
			"ping":    <-chans[0],
			"jitter":  <-chans[1],
			"monitor": <-chans[2],
		}
	} else {
		var (
			startTimestamp, endTimestamp int64
			err                          error
		)
		if startTimestampStr == "" {
			startTimestamp = int64(math.MaxInt64)
		} else {
			startTimestamp, err = strconv.ParseInt(startTimestampStr, 10, 64)
			if err != nil {
				logger.Terminal(fmt.Errorf("error in startTimestamp: %s", err.Error()).Error(), "p")
			}
		}
		if endTimestampStr == "" {
			endTimestamp = int64(math.MinInt64)
		} else {
			endTimestamp, err = strconv.ParseInt(endTimestampStr, 10, 64)
			if err != nil {
				logger.Terminal(fmt.Errorf("error in endTimestamp: %s", err.Error()).Error(), "p")
			}
		}

		go parallelQueryExec(matrix.PingChain.Path, querier.TypeRange, startTimestamp, endTimestamp, chans[0])
		go parallelQueryExec(matrix.JitterChain.Path, querier.TypeRange, startTimestamp, endTimestamp, chans[1])
		go parallelQueryExec(matrix.MonitorChain.Path, querier.TypeRange, startTimestamp, endTimestamp, chans[2])
		matrixResponse = map[string]querier.QueryResponse{
			"ping":    <-chans[0],
			"jitter":  <-chans[1],
			"monitor": <-chans[2],
		}
	}
	a.Data = matrixResponse
	a.send(w, a.marshalled())
}

// TSDBPathDetails responds with the path details that will be used for
// passing into the querier's timeSeriesPath.
func (a *API) TSDBPathDetails(w http.ResponseWriter, _ *http.Request) {
	var chainDetails []utils.ResponseTSDBChains
	for n, v := range *a.Matrix {
		chainDetails = append(chainDetails, utils.ResponseTSDBChains{
			Name: v.Domain,
			Path: utils.ChainPath{
				MatrixName: n,
				Ping:       trim(v.PingChain.Path),
				Jitter:     trim(v.JitterChain.Path),
				Fping:      trim(v.FPingChain.Path),
				Monitor:    trim(v.MonitorChain.Path),
			},
		})
	}
	a.Data = chainDetails
	a.send(w, a.marshalled())
}

// UpdateMonitoringServicesState starts the monitoring services on request from the API.
func (a *API) UpdateMonitoringServicesState(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != "start" && state != "stop" {
		panic("start-monitoring: invalid state received: " + state)
	}

	service := reflect.ValueOf(a.Services).Elem()
	sp, ok := service.FieldByName("Ping").Interface().(*ping.Ping)
	if !ok {
		panic("start-monitoring: ping not found")
	}
	if !sp.Iterate(state, false) {
		panic("start-monitoring: triggering monitoring: ping")
	}

	sj, ok := service.FieldByName("Jitter").Interface().(*jitter.Jitter)
	if !ok {
		panic("start-monitoring: jitter not found")
	}
	if !sj.Iterate(state, false) {
		panic("start-monitoring: triggering monitoring: jitter")
	}

	sm, ok := service.FieldByName("Monitor").Interface().(*monitor.Monitor)
	if !ok {
		panic("start-monitoring: monitor not found")
	}
	if !sm.Iterate(state, false) {
		panic("start-monitoring: triggering monitoring: monitor")
	}

	a.Data = true
	a.send(w, a.marshalled())
}

// GetMonitoringState returns the monitoring state.
func (a *API) GetMonitoringState(w http.ResponseWriter, r *http.Request) {
	services := a.config.Config.UtilsConf.ServicesSignal
	if services.Jitter != services.Ping && services.Jitter != services.ReqResDelayMonitoring {
		panic("get-monitoring-state: services state not aligned")
	}
	a.Data = services.Jitter
	a.send(w, a.marshalled())
}

func (a *API) marshalled() []byte {
	response := struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}{
		Status: a.ResponseStatus,
		Data:   a.Data,
	}
	js, err := json.Marshal(response)
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

func trim(s string) string {
	return s[0 : len(s)-5]
}
