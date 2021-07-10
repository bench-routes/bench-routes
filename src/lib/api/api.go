package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	router    *mux.Router
	reloadSig chan<- struct{}
}

func New(reloadSig chan<- struct{}) *API {
	tmp := &API{
		router:    mux.NewRouter(),
		reloadSig: reloadSig,
	}
	tmp.RegisterRoutes()
	return tmp
}

func (a *API) Router() *mux.Router {
	return a.router
}

func (a *API) RegisterRoutes() {
	a.router.HandleFunc("/api/v1/reload", a.Reload)
}

func (a *API) Reload(w http.ResponseWriter, r *http.Request) {
	a.reloadSig <- struct{}{}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true, Data: "reloaded success"}`))
}
