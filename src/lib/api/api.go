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
	log.Debug("msg", "Updated API domain map")
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
	a.send(w, http.StatusOK, `reload success`)
}

func (a *API) getMachines(w http.ResponseWriter, r *http.Request) {
	var res struct {
		Machines []string `json:"machines"`
	}
	for k := range a.domainMap {
		res.Machines = append(res.Machines, k)
	}
	a.send(w, http.StatusOK, res)
}

func (a *API) getDomainEntity(w http.ResponseWriter, r *http.Request) {
	domain_or_ip := r.URL.Query().Get("domain_or_ip")
	type data struct {
		ChainName  string `json:"chain_name"`
		EntityName string `json:"entity_name"`
		Status     string `json:"status"`
	}
	res := make([]data, 0)
	// append ping path if it exists.
	// ping path can be structured as <domain_or_ip> + "_ping" + <tsdb fileExtension>
	pingPath := "./storage/" + domain_or_ip + "_ping" + tsdb.FileExtension
	if ok := tsdb.VerifyChainPathExists(pingPath); ok {
		res = append(res, data{ChainName: pingPath, EntityName: "Ping", Status: "none"})
	} else {
		log.Warn("component", "api", "msg", fmt.Sprintf("ping: path not found at %s", pingPath))
	}

	// append jitter path if it exists.
	// jitter path can be structured as <domain_or_ip> + "_jitter" + <tsdb fileExtension>
	jitterPath := "./storage/" + domain_or_ip + "_jitter" + tsdb.FileExtension
	if ok := tsdb.VerifyChainPathExists(jitterPath); ok {
		res = append(res, data{ChainName: jitterPath, EntityName: "Jitter", Status: "none"})
	} else {
		log.Warn("component", "api", "msg", fmt.Sprintf("jitter: path not found at %s", jitterPath))
	}

	for _, api := range a.domainMap[domain_or_ip] {
		//rest routes are for monitoring
		// append monitor path if it exists.
		// monitor path can be structured as <domain_or_ip> + "_monitor" + <tsdb fileExtension>
		monitorPath := "./storage/" + api.Name + "_monitor" + tsdb.FileExtension
		if ok := tsdb.VerifyChainPathExists(monitorPath); ok {
			qry, err := querier.New(querier.TypeFirst, monitorPath, math.MinInt64, math.MaxInt64)
			if err != nil {
				a.send(w, http.StatusInternalServerError, err.Error())
				return
			}
			qryRes, err := qry.Exec()
			if err != nil {
				a.send(w, http.StatusInternalServerError, err.Error())
				return
			}
			var status string
			// We verify if status is 200, and only then ask the UI to show `UP` for the route.
			if qryRes.Values[0].Value.(evaluate.Response).Status != 200 {
				status = "false"
			} else {
				status = "true"
			}
			res = append(res, data{ChainName: monitorPath, EntityName: api.Route, Status: status})
		} else {
			log.Warn("component", "api", "msg", fmt.Sprintf("monitoring: path not found at %s", monitorPath))
		}
	}
	a.send(w, http.StatusOK, res)
}

func (a *API) queryEntity(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("name")
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	stepStr := r.URL.Query().Get("step")
	if path == "" {
		a.send(w, http.StatusBadRequest, fmt.Errorf("`name` query missing").Error())
		return
	}

	start, err := parseTime(startStr, "start")
	if err != nil {
		a.send(w, http.StatusBadRequest, err.Error())
		return
	}
	end, err := parseTime(endStr, "end")
	if err != nil {
		a.send(w, http.StatusBadRequest, err.Error())
		return
	}
	step, err := parseStep(stepStr)
	if err != nil {
		a.send(w, http.StatusBadRequest, err.Error())
		return
	}
	if ok := tsdb.VerifyChainPathExists(path); !ok {
		a.send(w, http.StatusBadRequest, fmt.Errorf("%s", "INVALID_PATH").Error())
		return
	}

	qry, err := querier.New(querier.TypeRange, path, start, end)
	if err != nil {
		a.send(w, http.StatusInternalServerError, err.Error())
		return
	}
	qryRes, err := qry.Exec()
	if err != nil {
		a.send(w, http.StatusInternalServerError, err.Error())
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
	a.send(w, http.StatusOK, res)
}

func (a *API) send(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	res := makeResponse(data, status)
	resp, err := json.Marshal(res)
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

func makeResponse(data interface{}, status int) *response {
	return &response{
		Status: status,
		Data:   data,
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
