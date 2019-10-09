package service

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/zairza-cetb/bench-routes/src/service/controllers"
)

var upgrader = websocket.Upgrader{
	// set buffer size to 3 mega-bytes
	ReadBufferSize:  3072,
	WriteBufferSize: 3072,
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Printf("ping\n")
	_, e := w.Write([]byte("pong"))
	if e != nil {
		panic(e)
	}
}

func sockets(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("error using upgrader %s\n", err)
	}
	log.Printf("connected from %s", r.URL)

	// capture client request for enabling series of responses unless its killed
	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("connection to client lost.\n%s\n", err)
			return
		}
		messageStr := string(message)
		log.Printf("type: %d\n message: %s \n", messageType, messageStr)

		switch messageStr {
		case "force-start-ping":
			// true if success else false
			e := ws.WriteMessage(1, []byte(strconv.FormatBool(controllers.PingController("start"))))
			if e != nil {
				panic(e)
			}
		case "force-stop-ping":
			// true if success else false
			e := ws.WriteMessage(1, []byte(strconv.FormatBool(controllers.PingController("stop"))))
			if e != nil {
				panic(e)
			}
		case "force-start-jitter":
			//true if success else false
			e := ws.WriteMessage(1, []byte(strconv.FormatBool(controllers.JitterController("start"))))
			if e != nil {
				panic(e)
			}
		case "force-stop-jitter":
			//true if success else false
			e := ws.WriteMessage(1, []byte(strconv.FormatBool(controllers.JitterController("stop"))))
			if e != nil {
				panic(e)
			}
		case "force-start-req-res-monitoring":
			// true if success else false
			e := ws.WriteMessage(1, []byte(strconv.FormatBool(controllers.ReqResMonitoringController("start"))))
			if e != nil {
				panic(e)
			}
		case "force-stop-req-res-monitoring":
			// true if success else false
			e := ws.WriteMessage(1, []byte(strconv.FormatBool(controllers.ReqResMonitoringController("stop"))))
			if e != nil {
				panic(e)
			}
		}
	}

}

//Service initiates services and avail different routes for bench-routes
func service(port string) {
	http.HandleFunc("/", home)
	http.HandleFunc("/websocket", sockets)

	log.Fatal(http.ListenAndServe(port, nil))
}

// Runner serves as a default runner for bench-routes
func Runner(port string) {
	log.Println("Starting bench-routes service ...")
	service(port)
}
