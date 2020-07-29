package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/pprof"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	config "github.com/zairza-cetb/bench-routes/src/lib/config"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/jitter"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/monitor"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/ping"
	"github.com/zairza-cetb/bench-routes/src/lib/request"
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
	ResponseStatus      string
	Data, Services      interface{}
	Matrix              *map[string]*utils.BRMatrix
	configurationPath   string
	config              *config.Config
	reloadConfigURLs    chan struct{}
	receiveFinishSignal chan struct{}
	mux                 sync.RWMutex
}

type inputRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Params  map[string]string `json:"params"`
	Headers map[string]string `json:"headers"`
	Body    map[string]string `json:"body"`
	Labels  []string          `json:"labels"`
}

// New returns the API type for implementing the API interface.
func New(
	matrix *map[string]*utils.BRMatrix,
	config *config.Config,
	configPath string,
	services interface{},
	reload,
	done chan struct{},
) *API {
	return &API{
		Matrix:              matrix,
		Services:            services,
		configurationPath:   configPath,
		config:              config,
		reloadConfigURLs:    reload,
		receiveFinishSignal: done,
	}
}

// Register registers the routes with the mux router.
func (a *API) Register(router *mux.Router) {
	// Static servings.
	{
		router.Handle("/", http.FileServer(http.Dir(uiPathV11)))
		router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(uiPathV11+"assets/"))))
		router.PathPrefix("/manifest.json").Handler(http.StripPrefix("/manifest.json", http.FileServer(http.Dir(uiPathV11+"/manifest.json"))))
		router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(uiPathV11+"static/"))))
	}
	// Profiling routes.
	{
		router.HandleFunc("/debug/pprof/", pprof.Index)
		router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		router.Handle("/debug/pprof/block", pprof.Handler("block"))
		router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	}
	// Prometheus metrics.
	{
		router.Handle("/metrics", promhttp.Handler())
	}
	// API endpoints.
	{
		router.HandleFunc("/add-route", a.AddRouteToMonitoring)
		router.HandleFunc("/br-live-check", a.Home)
		router.HandleFunc("/config/update-interval", a.ModifyIntervalDuration)
		router.HandleFunc("/delete-route", a.DeleteConfigRoutes)
		router.HandleFunc("/get-monitoring-services-state", a.GetMonitoringState)
		router.HandleFunc("/get-config-intervals", a.GetConfigIntervals)
		router.HandleFunc("/get-config-routes", a.GetConfigRoutes)
		router.HandleFunc("/get-labels", a.GetLabels)
		router.HandleFunc("/get-route-time-series", a.TSDBPathDetails)
		router.HandleFunc("/query-matrix", a.SendMatrix)
		router.HandleFunc("/query", a.Query)
		router.HandleFunc("/quick-input", a.QuickTestInput)
		router.HandleFunc("/routes-summary", a.RoutesSummary)
		router.HandleFunc("/service-state", a.ServiceState)
		router.HandleFunc("/test", a.TestTemplate)
		router.HandleFunc("/update-monitoring-services-state", a.UpdateMonitoringServicesState)
		router.HandleFunc("/update-route", a.UpdateRoute)
	}
}

// Home handles the requests for the home page.
func (a *API) Home(w http.ResponseWriter, r *http.Request) {
	msg := "ping from " + r.RemoteAddr + ", sent pong in monitor"
	log.Infoln(msg)
	a.Data = msg
	a.send(w, a.marshalled())
}

// UIv1 serves the v1.0 version of user-interface of bench-routes.
// ui-builds/v1.0 is served through this.
func (a *API) UIv1(_ http.ResponseWriter, _ *http.Request) {
	http.FileServer(http.Dir(uiPathV1))
}

// TestTemplate handles template related request for testing purposes.
func (a *API) TestTemplate(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, testFilesDir+"bench-routes-socket-tester.html")
}

// ServiceState handles requests related to the state of the Services in
// the application.
func (a *API) ServiceState(w http.ResponseWriter, _ *http.Request) {
	p := config.New(a.configurationPath)
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
func (a *API) RoutesSummary(w http.ResponseWriter, _ *http.Request) {
	p := config.New(a.configurationPath)
	p.Refresh()

	var servicesRoutes, monitoringRoutes []string
	for _, r := range p.Config.Routes {
		servicesRoutes = append(servicesRoutes, r.URL)
		monitoringRoutes = append(monitoringRoutes, r.Method+": "+r.URL)
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
			log.Warnln(fmt.Errorf("error in startTimestamp: %s", err.Error()).Error())
		}
	}

	endTimestampStr := r.URL.Query().Get("endTimestamp")
	if endTimestampStr == "" {
		endTimestamp = int64(math.MinInt64)
	} else {
		endTimestamp, err = strconv.ParseInt(endTimestampStr, 10, 64)
		if err != nil {
			log.Warnln(fmt.Errorf("error in endTimestamp: %s", err.Error()).Error())
		}
	}

	// condition: only for bench-routes as per the design
	//
	// path should be in syntax: <DBPath>/<ofType>/chunk_<middle>_<url>.json -> non-system metric
	// %s/system.json -> system metric

	// verify if chain path exists
	timeSeriesPath = timeSeriesPath + tsdb.FileExtension
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
	routeHashMatrix := r.URL.Query().Get("routeNameMatrix")
	if _, ok := (*a.Matrix)[routeHashMatrix]; !ok {
		a.Data = "ROUTE_NAME_AKA_INSTANCE_KEY_NOT_IN_MATRIX"
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
	matrix := (*a.Matrix)[routeHashMatrix]
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
				log.Warnln(fmt.Errorf("error in startTimestamp: %s", err.Error()).Error())
			}
		}
		if endTimestampStr == "" {
			endTimestamp = int64(math.MinInt64)
		} else {
			endTimestamp, err = strconv.ParseInt(endTimestampStr, 10, 64)
			if err != nil {
				log.Warnln(fmt.Errorf("error in endTimestamp: %s", err.Error()).Error())
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

// QuickTestInput tests the API input from the /quick-input route page
// of the react-UI.
func (a *API) QuickTestInput(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	var (
		t       inputRequest
		decoder = json.NewDecoder(r.Body)
	)
	if err := decoder.Decode(&t); err != nil {
		panic(err)
	}
	req := request.New(t.URL, t.Headers, t.Params, t.Body, t.Labels)
	response := make(chan request.ResponseWrapper)
	switch t.Method {
	case "GET":
		go req.Send(request.GET, response)
	case "POST":
		go req.Send(request.POST, response)
	default:
		fmt.Printf("invalid request method: %s\n", t.Method)
	}
	a.Data = <-response
	a.send(w, a.marshalled())
}

// AddRouteToMonitoring adds a new route to the config.
func (a *API) AddRouteToMonitoring(w http.ResponseWriter, r *http.Request) {
	var (
		t       inputRequest
		decoder = json.NewDecoder(r.Body)
	)
	if err := decoder.Decode(&t); err != nil {
		panic(err)
	}
	requestInstance := request.New(t.URL, t.Headers, t.Params, t.Body, t.Labels)
	a.config.AddRoute(
		config.GetNewRouteType(
			t.Method,
			t.URL,
			requestInstance.GetHeadersConfigFormatted(),
			requestInstance.GetParamsConfigFormatted(),
			requestInstance.GetBodyConfigFormatted(),
			t.Labels,
		),
	)
	a.reloadConfigURLs <- struct{}{}
	a.Data = "success"
	a.send(w, a.marshalled())
	<-a.receiveFinishSignal
}

// TSDBPathDetails responds with the path details that will be used for
// passing into the querier's timeSeriesPath.
func (a *API) TSDBPathDetails(w http.ResponseWriter, _ *http.Request) {
	var chainDetails []utils.ResponseTSDBChains
	for hash, v := range *a.Matrix {
		chainDetails = append(chainDetails, utils.ResponseTSDBChains{
			Name: v.FullURL,
			Path: utils.ChainPath{
				InstanceKey: hash,
				Ping:        trim(v.PingChain.Path),
				Jitter:      trim(v.JitterChain.Path),
				Fping:       trim(v.FPingChain.Path),
				Monitor:     trim(v.MonitorChain.Path),
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
		log.Warnln("start-monitoring: invalid state received: " + state)
		a.Data = "INVALID_STATE"
	} else {
		a.Data = true
	}
	service := reflect.ValueOf(a.Services).Elem()
	sp, ok := service.FieldByName("Ping").Interface().(*ping.Ping)
	if !ok {
		panic("start-monitoring: ping not found")
	}

	sp.Iterate(state, false)

	sj, ok := service.FieldByName("Jitter").Interface().(*jitter.Jitter)
	if !ok {
		panic("start-monitoring: jitter not found")
	}
	sj.Iterate(state, false)

	sm, ok := service.FieldByName("Monitor").Interface().(*monitor.Monitor)
	if !ok {
		panic("start-monitoring: monitor not found")
	}
	sm.Iterate(state, false)
	a.send(w, a.marshalled())
}

// GetMonitoringState returns the monitoring state.
func (a *API) GetMonitoringState(w http.ResponseWriter, _ *http.Request) {
	service := reflect.ValueOf(a.Services).Elem()
	sp, ok := service.FieldByName("Ping").Interface().(*ping.Ping)
	if !ok {
		panic("start-monitoring: ping not found")
	}

	sj, ok := service.FieldByName("Jitter").Interface().(*jitter.Jitter)
	if !ok {
		panic("start-monitoring: jitter not found")
	}

	sm, ok := service.FieldByName("Monitor").Interface().(*monitor.Monitor)
	if !ok {
		panic("start-monitoring: monitor not found")
	}

	if sp.IsActive() != sj.IsActive() || sm.IsActive() != sp.IsActive() || sm.IsActive() != sj.IsActive() {
		panic("states not aligned")
	}

	if sp.IsActive() {
		a.Data = "active"
	} else {
		a.Data = "passive"
	}
	a.send(w, a.marshalled())
}

// GetConfigIntervals gets the config file data for the config screen.
func (a *API) GetConfigIntervals(w http.ResponseWriter, _ *http.Request) {
	a.Data = a.config.Config.Interval
	a.send(w, a.marshalled())
}

// GetConfigRoutes gets the config file data for the config screen.
func (a *API) GetConfigRoutes(w http.ResponseWriter, _ *http.Request) {
	a.Data = a.config.Config.Routes
	a.send(w, a.marshalled())
}

// ModifyIntervalDuration modifies a specific interval duration in the config file.
func (a *API) ModifyIntervalDuration(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	var req struct {
		IntervalName string `json:"intervalName"`
		Value        string `json:"value"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		panic(err)
	}
	if num, err := strconv.Atoi(req.Value); err == nil {
		n := int64(num)
		for i, interval := range a.config.Config.Interval {
			if interval.Test == req.IntervalName {
				a.config.Config.Interval[i].Duration = &n
			}
		}
		_, err := a.config.Write()
		if err == nil {
			a.ResponseStatus = "200"
			a.Data = a.config
		} else {
			a.ResponseStatus = "400"
			a.Data = "Could not modify the config file"
		}
	} else {
		// The string is not an integer.
		a.ResponseStatus = "400"
		a.Data = "The string passed is not an integer"
	}
	a.send(w, a.marshalled())
}

// UpdateRoute updates a route in the local config.
func (a *API) UpdateRoute(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	var (
		req struct {
			Method        string            `json:"method"`
			URL           string            `json:"url"`
			Params        map[string]string `json:"params"`
			Headers       map[string]string `json:"headers"`
			Body          map[string]string `json:"body"`
			OriginalRoute string            `json:"orgRoute"`
			Labels        []string          `json:"labels"`
		}
		decoder = json.NewDecoder(r.Body)
	)
	if err := decoder.Decode(&req); err != nil {
		panic(err)
	}
	requestInstance := request.New(req.URL, req.Headers, req.Params, req.Body, req.Labels)
	for i, route := range a.config.Config.Routes {
		if route.URL == req.OriginalRoute && route.Method == req.Method {
			a.config.Config.Routes = append(a.config.Config.Routes[:i], a.config.Config.Routes[i+1:]...)
			a.config.Config.Routes = append(a.config.Config.Routes, config.GetNewRouteType(
				req.Method,
				req.URL,
				requestInstance.GetHeadersConfigFormatted(),
				requestInstance.GetParamsConfigFormatted(),
				requestInstance.GetBodyConfigFormatted(),
				req.Labels,
			))
			break
		}
	}
	if _, err := a.config.Write(); err != nil {
		a.ResponseStatus = http.StatusText(400)
	}
	a.reloadConfigURLs <- struct{}{}
	<-a.receiveFinishSignal
	a.ResponseStatus = http.StatusText(200)
	a.Data = a.config.Config.Routes
	a.send(w, a.marshalled())
}

// DeleteConfigRoutes removes a route from the config screen.
func (a *API) DeleteConfigRoutes(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	var req struct {
		ActualRoute string `json:"actualRoute"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		panic(err)
	}
	for i, route := range a.config.Config.Routes {
		if route.URL == req.ActualRoute {
			a.mux.Lock()
			a.config.Config.Routes = append(a.config.Config.Routes[:i], a.config.Config.Routes[i+1:]...)
			a.mux.Unlock()
			break
		}
	}
	if _, err := a.config.Write(); err != nil {
		a.ResponseStatus = http.StatusText(400)
	}
	a.ResponseStatus = http.StatusText(200)
	a.Data = a.config.Config.Routes
	a.send(w, a.marshalled())
}

// GetLabels gets route labels from the config file
func (a *API) GetLabels(w http.ResponseWriter, r *http.Request) {
	var (
		uniqueLabels []string
		m            = make(map[string]bool)
	)
	for _, route := range a.config.Config.Routes {
		for _, label := range route.Labels {
			if _, value := m[label]; !value {
				m[label] = true
				uniqueLabels = append(uniqueLabels, label)
			}
		}
	}
	a.ResponseStatus = http.StatusText(200)
	a.Data = uniqueLabels
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
