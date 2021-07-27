package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/bench-routes/bench-routes/src/lib/log"
	"github.com/bench-routes/bench-routes/src/lib/modules/evaluate"
	tsdb "github.com/bench-routes/bench-routes/tsdb/file"
	"github.com/bench-routes/bench-routes/tsdb/file/querier"
	"github.com/gorilla/mux"
)

type API struct {
	router    *mux.Router
	reloadSig chan<- struct{}
	mutex     sync.RWMutex
	domainMap map[string][]config.API
}

type response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type errResponse struct {
	Error string `json:"error"`
}

func New(reloadSig chan<- struct{}, conf *config.Config) *API {
	tmp := &API{
		router:    mux.NewRouter(),
		reloadSig: reloadSig,
	}
	tmp.UpdateConf(conf)
	tmp.RegisterRoutes()
	return tmp
}

func (a *API) Router() *mux.Router {
	return a.router
}

func (a *API) UpdateConf(conf *config.Config) {
	d := make(map[string][]config.API)
	for _, api := range conf.APIs {
		d[api.Domain] = append(d[api.Domain], api)
	}
	a.mutex.Lock()
	a.domainMap = d
	log.Info("msg", "Updated API domain map")
	a.mutex.Unlock()
}

func (a *API) RegisterRoutes() {
	a.router.HandleFunc("/api/v1/reload", a.Reload)
	a.router.HandleFunc("/api/v1/get-machines", a.getMachines)
	a.router.HandleFunc("/api/v1/get-domain-entities", a.getDomainEntity)
	a.router.HandleFunc("/api/v1/query-entity", a.queryEntity)
}

func (a *API) Reload(w http.ResponseWriter, r *http.Request) {
	a.reloadSig <- struct{}{}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true, "data": "reload success"}`))
}

func (a *API) getMachines(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Machines []string `json:"machines"`
	}
	var res data
	for k := range a.domainMap {
		res.Machines = append(res.Machines, k)
	}
	a.send(w, makeResponse(res), http.StatusOK)
}

func (a *API) getDomainEntity(w http.ResponseWriter, r *http.Request) {
	domain_or_ip := r.URL.Query().Get("domain_or_ip")
	type data struct {
		ChainName  string `json:"chain_name"`
		EntityName string `json:"entity_name"`
		Status     string `json:"status"`
	}
	res := []data{}
	// append ping path if it exists.
	// ping path can be structured as <domain_or_ip> + "_ping" + <tsdb fileExtension>
	pingPath := "./storage/" + domain_or_ip + "_ping" + tsdb.FileExtension
	if ok := tsdb.VerifyChainPathExists(pingPath); ok {
		res = append(res, data{ChainName: pingPath, EntityName: "Ping", Status: "none"})
	}

	// append jitter path if it exists.
	// jitter path can be structured as <domain_or_ip> + "_jitter" + <tsdb fileExtension>
	jitterPath := "./storage/" + domain_or_ip + "_jitter" + tsdb.FileExtension
	if ok := tsdb.VerifyChainPathExists(jitterPath); ok {
		res = append(res, data{ChainName: jitterPath, EntityName: "Jitter", Status: "none"})
	}

	for _, api := range a.domainMap[domain_or_ip] {
		//rest routes are for monitoring
		// append monitor path if it exists.
		// monitor path can be structured as <domain_or_ip> + "_monitor" + <tsdb fileExtension>
		monitorPath := "./storage/" + api.Name + "_monitor" + tsdb.FileExtension
		if ok := tsdb.VerifyChainPathExists(monitorPath); ok {
			qry, err := querier.New(querier.TypeFirst, monitorPath, math.MinInt64, math.MaxInt64)
			if err != nil {
				a.send(w, makeError(err), http.StatusBadRequest)
				return
			}
			qryRes, err := qry.Exec()
			if err != nil {
				a.send(w, makeError(err), http.StatusBadRequest)
				return
			}
			var status string
			if qryRes.Values[0].Value.(evaluate.Response).Status != 200 {
				status = "false"
			} else {
				status = "true"
			}
			res = append(res, data{ChainName: monitorPath, EntityName: api.Route, Status: status})
		}
	}
	a.send(w, makeResponse(res), http.StatusOK)
}

func (a *API) queryEntity(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("name")
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	stepStr := r.URL.Query().Get("step")
	if path == "" {
		a.send(w, makeError(fmt.Errorf("`name` query missing")), http.StatusBadRequest)
		return
	}

	start, err := parseTime(startStr, "start")
	if err != nil {
		a.send(w, makeError(err), http.StatusBadRequest)
		return
	}
	end, err := parseTime(endStr, "end")
	if err != nil {
		a.send(w, makeError(err), http.StatusBadRequest)
		return
	}
	step, err := parseStep(stepStr)
	if err != nil {
		a.send(w, makeError(err), http.StatusBadRequest)
		return
	}
	if ok := tsdb.VerifyChainPathExists(path); !ok {
		a.send(w, makeError(fmt.Errorf("%s", "INVALID_PATH")), http.StatusBadRequest)
		return
	}

	qry, err := querier.New(querier.TypeRange, path, start, end)
	if err != nil {
		a.send(w, makeError(err), http.StatusBadRequest)
		return
	}
	qryRes, err := qry.Exec()
	if err != nil {
		a.send(w, makeError(err), http.StatusBadRequest)
		return
	}

	type block struct {
		Timestamp int64       `json:"timestamp"`
		Value     interface{} `json:"value"`
	}

	type data struct {
		Unit string  `json:"unit"`
		Data []block `json:"data"`
	}
	var res data
	for i := 0; i < len(qryRes.Values); i += (1 + step) {
		switch qryRes.Type {
		case "jitter":
			res.Unit = "ms"
			res.Data = append(res.Data, block{
				Timestamp: qryRes.Values[i].NormalizedTime / int64(time.Millisecond),
				Value:     qryRes.Values[i].Value.(evaluate.Jitter).Value.Milliseconds(),
			})
		case "ping":
			res.Unit = "ms"
			res.Data = append(res.Data, block{
				Timestamp: qryRes.Values[i].NormalizedTime / int64(time.Millisecond),
				Value:     qryRes.Values[i].Value.(evaluate.Ping).Mean.Milliseconds(),
			})
		case "monitoring":
			res.Unit = "ms"
			res.Data = append(res.Data, block{
				Timestamp: qryRes.Values[i].NormalizedTime / int64(time.Millisecond),
				Value:     qryRes.Values[i].Value.(evaluate.Response).Delay.Milliseconds(),
			})
		}
	}
	a.send(w, makeResponse(res), http.StatusOK)
}

func (a *API) send(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	resp, err := json.Marshal(data)
	if err != nil {
		fmt.Println(fmt.Errorf("error marshalling : %w", err))
		return
	}
	w.Write([]byte(resp))
}

func parseTime(t string, typ string) (int64, error) {
	if t == "" && typ == "end" {
		return int64(math.MaxInt64), nil
	}
	if t == "" && typ == "start" {
		return int64(math.MinInt64), nil
	}
	time, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return 0, fmt.Errorf("error in timestamp: %s", err.Error())
	}
	return time.UnixNano(), nil
}

func makeResponse(i interface{}) *response {
	return &response{
		Status: http.StatusOK,
		Data:   i,
	}
}

func makeError(err error) *errResponse {
	fmt.Println(err.Error())
	return &errResponse{
		Error: err.Error(),
	}
}

func parseStep(stepStr string) (int, error) {
	if stepStr == "" {
		return 0, nil
	}
	step, err := strconv.Atoi(stepStr)
	if err != nil {
		return 0, fmt.Errorf("step parsing error: %s", err.Error())
	}
	return step, nil
}
