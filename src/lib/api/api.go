package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/bench-routes/bench-routes/src/lib/log"
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
	a.router.HandleFunc("/get-machines", a.getMachines)
	a.router.HandleFunc("/get-domain-entities", a.getDomainEntity)
	// a.router.HandleFunc("/query-entity", a.Reload)
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
	type response struct {
		Name   string `json:"name"`
		Route  string `json:"route"`
		Status bool   `json:"status"`
	}
	res := []response{}
	domain_or_ip := r.URL.Query().Get("domain_or_ip")
	for _, api := range a.domainMap[domain_or_ip] {
		res = append(res, response{Name: api.Name, Route: api.Route, Status: true})
	}
	a.send(w, res, http.StatusOK)
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
