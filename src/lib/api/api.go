package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/bench-routes/bench-routes/src/lib/log"
	tsdb "github.com/bench-routes/bench-routes/tsdb/file"
	"github.com/bench-routes/bench-routes/tsdb/querier"
	"github.com/gorilla/mux"
)

type API struct {
	router    *mux.Router
	reloadSig chan<- struct{}
	mutex     sync.RWMutex
	domainMap map[string][]config.API
}

func New(reloadSig chan<- struct{}, conf *config.Config) *API {
	tmp := &API{
		router:    mux.NewRouter(),
		reloadSig: reloadSig,
		domainMap: make(map[string][]config.API),
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
	// a.router.HandleFunc("/query-entity", a.queryEntity)
}

func (a *API) Reload(w http.ResponseWriter, r *http.Request) {
	a.reloadSig <- struct{}{}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true, "Data": "reloaded success"}`))
}

func (a *API) getMachines(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Machines []string `json:"machines"`
	}
	var res response
	for k := range a.domainMap {
		res.Machines = append(res.Machines, k)
	}
	a.send(w, res, http.StatusOK)
}

func (a *API) getDomainEntity(w http.ResponseWriter, r *http.Request) {
	domain_or_ip := r.URL.Query().Get("domain_or_ip")
	type response struct {
		ChainName  string `json:"name"`
		EntityName string `json:"route"`
		Status     string `json:"status"`
	}
	res := []response{}
	// append ping path if it exists.
	// ping path can be structured as <domain_or_ip> + "_ping" + <tsdb fileExtension>
	pingPath := "./storage/" + domain_or_ip + "_ping" + tsdb.FileExtension
	if ok := tsdb.VerifyChainPathExists(pingPath); ok {
		res = append(res, response{ChainName: pingPath, EntityName: "Ping", Status: "none"})
	}

	// append jitter path if it exists.
	// jitter path can be structured as <domain_or_ip> + "_jitter" + <tsdb fileExtension>
	jitterPath := "./storage/" + domain_or_ip + "_jitter" + tsdb.FileExtension
	if ok := tsdb.VerifyChainPathExists(jitterPath); ok {
		res = append(res, response{ChainName: jitterPath, EntityName: "Jitter", Status: "none"})
	}

	for _, api := range a.domainMap[domain_or_ip] {
		//rest routes are for monitoring
		// append monitor path if it exists.
		// monitor path can be structured as <domain_or_ip> + "_monitor" + <tsdb fileExtension>
		monitorPath := "./storage/" + api.Name + "_monitor" + tsdb.FileExtension
		if ok := tsdb.VerifyChainPathExists(monitorPath); ok {
			res = append(res, response{ChainName: monitorPath, EntityName: api.Route, Status: "true"})
		}
	}
	a.send(w, res, http.StatusOK)
}

func (a *API) queryEntity(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("name")
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	start, err := parseTime(startStr)
	if err != nil {
		a.send(w, err.Error(), http.StatusBadRequest)
		return
	}
	end, err := parseTime(endStr)
	if err != nil {
		a.send(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = path + tsdb.FileExtension
	if ok := tsdb.VerifyChainPathExists(path); !ok {
		a.send(w, "INVALID_PATH", http.StatusBadRequest)
		return
	}

	qry := querier.New(path, "", querier.TypeRange)
	query := qry.QueryBuilder()
	query.SetRange(start, end)
	a.send(w, query.ExecWithoutEncode(), http.StatusOK)
}

func parseTime(t string) (int64, error) {
	if t == "" {
		return int64(math.MaxInt64), nil
	}
	time, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return 0, fmt.Errorf("error in timestamp: %s", err.Error())
	}
	return time.UnixNano(), nil
}

func (a *API) send(w http.ResponseWriter, response interface{}, status int) {
	w.WriteHeader(status)
	resp, err := json.Marshal(response)
	if err != nil {
		fmt.Println(fmt.Errorf("error marshalling : %w", err))
		return
	}
	w.Write([]byte(resp))
}
