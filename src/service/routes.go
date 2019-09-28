package service

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/zairza-cetb/bench-routes/src/service/controllers"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  3072,
	WriteBufferSize: 3072,
}

func init() {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.Println("init started")
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Printf("ping\n")
	w.Write([]byte("pong"))
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
			log.Fatalf("error rendering main socket %s\n", err)
			return
		}
		messageStr := string(message)
		log.Printf("type: %d\n message: %s \n", messageType, messageStr)

		switch messageStr {
		case "force-start-ping":
			controllers
			break
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
