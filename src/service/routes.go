package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zairza-cetb/bench-routes/src/service/controllers"
)

func init() {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.Println("init started")
}

func error404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, `404 page not found`)
}

//Service initiates services and avail different routes for bench-routes
func Service(port string) {

	r := mux.NewRouter()
	fmt.Println("Running on port", port)
	r.HandleFunc("/login", controllers.LoginHandler)

	r.NotFoundHandler = http.HandlerFunc(error404)

	log.Fatal(http.ListenAndServe(port, r))
}

// Runner serves as a default runner for bench-routes
func Runner(port string) {
	log.Println("Init Started")
	Service(port)
}
