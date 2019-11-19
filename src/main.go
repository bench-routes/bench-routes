package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/src/lib/utils/parser"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

var (
	port     = ":9090"
	upgrader = websocket.Upgrader{
		// set buffer to 4 mega-bytes size
		ReadBufferSize:  4048,
		WriteBufferSize: 4048,
	}
	configuration parser.YAMLBenchRoutesType
)

const (
	logFilePrefix = "bench-route-"
	logDirectory  = "br-logs"
	testFilesDir  = "test-files/"
)

func init() {
	setupLogger()
	log.Printf("initializing bench-routes ...")

	// load configuration file
	configuration.Address = utils.ConfigurationFilePath
	configuration = *configuration.Load()

	var ConfigURLs []string

	//function to initialize the service state before starting
	initializeState(&configuration)

	// Load and build TSDB chain
	// searching for unique URLs
	for _, r := range configuration.Config.Routes {
		found := false
		for _, i := range ConfigURLs {
			if i == r.URL {
				found = true
				break
			}
		}
		if !found {
			filters.HTTPPingFilter(&r.URL)
			ConfigURLs = append(ConfigURLs, r.URL)
			utils.PingDBNames[r.URL] = utils.GetHash(r.URL)
			utils.FloodPingDBNames[r.URL] = utils.GetHash(r.URL)
		}
	}
	// forming ping chain
	for _, v := range ConfigURLs {
		path := utils.PathPing + "/" + "chunk_ping_" + v + ".json"
		inst := &tsdb.Chain{
			Path:           path,
			Chain:          []tsdb.Block{},
			LengthElements: 0,
			Size:           0,
		}
		inst.Init().Save()
		// Initiate the chain
		utils.GlobalPingChain = append(utils.GlobalPingChain, inst)
	}

	// forming flood-ping chain
	for _, v := range ConfigURLs {
		path := utils.PathFloodPing + "/" + "chunk_flood_ping_" + v + ".json"
		inst := &tsdb.Chain{
			Path:           path,
			Chain:          []tsdb.Block{},
			LengthElements: 0,
			Size:           0,
		}
		inst.Init().Save()
		// Initiate the chain
		utils.GlobalFloodPingChain = append(utils.GlobalFloodPingChain, inst)
	}

	for _, v := range ConfigURLs {
		path := utils.PathJitter + "/" + "chunk_jitter_" + v + ".json"
		inst := &tsdb.Chain{
			Path:           path,
			Chain:          []tsdb.Block{},
			LengthElements: 0,
			Size:           0,
		}
		inst.Init().Save()
		// Initiate the chain
		utils.GlobalChain = append(utils.GlobalChain, inst)
	}

	// forming req-res-delay chain
	for _, route := range configuration.Config.Routes {
		path := utils.PathReqResDelayMonitoring + "/" + "chunk_req_res_" + filters.RouteDestroyer(route.URL)
		// Create sample chains to init in each TSDB file
		resp := &tsdb.Chain{
			Path:           path + ".json",
			Chain:          []tsdb.Block{},
			LengthElements: 0,
			Size:           0,
		}
		resp.Init().Save()
		utils.GlobalReqResDelChain = append(utils.GlobalReqResDelChain, resp)
	}

	// keep the below line to the end of file so that we ensure that we give a confirmation message only when all the
	// required resources for the application is up and healthy
	log.Printf("Bench-routes is up and running\n")
}

func main() {

	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("ping from %s, sent pong in response\n", r.RemoteAddr)
	})
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, testFilesDir+"bench-routes-socket-tester.html")
	})
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatalf("error using upgrader %s\n", err)
		}

		// capture client request for enabling series of responses unless its killed
		for {
			messageType, message, err := ws.ReadMessage()
			if err != nil {
				log.Printf("connection to client lost.\n%s\n", err)
				return
			}

			// In order to support compound signals, we aim to format the messages as:
			// <signal-name> <[optional] data>
			// The first param refers to the signal for the operation to be carried out.
			// The second param [optional] is a JSON object (stringified) which would be used for
			// general communication with the UI.
			// For example: > Qping-route {"url": "https://www.google.co.in"}		(compound signal)
			// 				> force-start-ping										(simple signal)
			inStream := strings.Split(string(message), " ")

			sig := inStream[0] // Signal
			log.Printf("type: %d\n message: %s \n", messageType, sig)
			// generate appropriate signals from incoming messages
			switch sig {
			// ping
			case "force-start-ping":
				// true if success else false
				if e := ws.WriteMessage(1, []byte(strconv.FormatBool(HandlerPingGeneral("start")))); e != nil {
					panic(e)
				}
			case "force-stop-ping":
				e := ws.WriteMessage(1, []byte(strconv.FormatBool(HandlerPingGeneral("stop"))))
				if e != nil {
					panic(e)
				}

				// flood-ping
			case "force-start-flood-ping":
				if e := ws.WriteMessage(1, []byte(strconv.FormatBool(HandlerFloodPingGeneral("start")))); e != nil {
					panic(e)
				}
			case "force-stop-flood-ping":
				if e := ws.WriteMessage(1, []byte(strconv.FormatBool(HandlerFloodPingGeneral("stop")))); e != nil {
					panic(e)
				}

				// jitter
			case "force-start-jitter":
				if e := ws.WriteMessage(1, []byte(strconv.FormatBool(HandlerJitterGeneral("start")))); e != nil {
					panic(e)
				}
			case "force-stop-jitter":
				if e := ws.WriteMessage(1, []byte(strconv.FormatBool(HandlerJitterGeneral("stop")))); e != nil {
					panic(e)
				}

				// request-response-monitoring
			case "force-start-req-res-monitoring":
				if e := ws.WriteMessage(1, []byte(strconv.FormatBool(HandleReqResGeneral("start")))); e != nil {
					panic(e)
				}
			case "force-stop-req-res-monitoring":
				if e := ws.WriteMessage(1, []byte(strconv.FormatBool(HandleReqResGeneral("stop")))); e != nil {
					panic(e)
				}

				// Get config routes details
			case "route-details":
				m := configuration.Config.Routes
				if e := ws.WriteMessage(1, filters.RouteYAMLtoJSONParser(m)); e != nil {
					panic(e)
				}

				// Queries
			case "Qping-route":
				compMessage := getMessageFromCompoundSignal(inStream[1:])
				inst := qPingRoute{}
				if e := json.Unmarshal(compMessage, &inst); e != nil {
					panic(e)
				}

				url := inst.URL
				ql := getQuerier(ws, "ping", url, "", "")
				go ql.FetchAllSeries()

			case "Qjitter-route":
				compMessage := getMessageFromCompoundSignal(inStream[1:])
				inst := qJitterRoute{}
				if e := json.Unmarshal(compMessage, &inst); e != nil {
					panic(e)
				}

				url := inst.URL
				ql := getQuerier(ws, "jitter", url, "", "")
				go ql.FetchAllSeries()

			case "Qflood-ping-route":
				compMessage := getMessageFromCompoundSignal(inStream[1:])
				inst := qFloodPingRoute{}
				if e := json.Unmarshal(compMessage, &inst); e != nil {
					panic(e)
				}

				url := inst.URL
				ql := getQuerier(ws, "flood-ping", url, "", "")
				go ql.FetchAllSeries()

			// Querrier signal for Request-response delay
			case "Qrequest-response-delay":
				compMessage := getMessageFromCompoundSignal(inStream[1:])
				inst := qReqResDelayRoute{}
				if e := json.Unmarshal(compMessage, &inst); e != nil {
					panic(e)
				}

				url := inst.URL
				method := inst.Method
				// Gets the Querrier for request-response delay
				// TODO: Send the method along with URL
				ql := getQuerier(ws, "req-res-delay", url, method, "_delay")
				go ql.FetchAllSeries()
			}
		}
	})

	// launch service
	log.Fatal(http.ListenAndServe(port, nil))

}

func setupLogger() {
	currTime := time.Now()
	currFileName := fmt.Sprint(logFilePrefix, currTime.Format("2006-01-02#15:04:05"), ".log")
	user, err := user.Current()
	if err != nil {
		fmt.Printf("cannot access current user data\n")
		return
	}

	homePath := user.HomeDir
	logDirectoryPath := homePath + "/" + logDirectory
	err = os.MkdirAll(logDirectoryPath, 0755)
	if err != nil {
		fmt.Printf("error creating log directory : %s\n", logDirectoryPath)
		return
	}
	logFilePath := logDirectoryPath + "/" + currFileName
	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0444)
	if err != nil {
		fmt.Printf("error opening log file : %s\n", logFilePath)
		return
	}
	writer := io.MultiWriter(os.Stdout, file)
	log.SetOutput(writer)
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func getQuerier(conn *websocket.Conn, serviceName, d, method, suff string) (inst tsdb.BRQuerier) {
	inst = tsdb.BRQuerier{
		ServiceName: serviceName,
		Route:       tsdb.BQRoute{DomainIP: d, Method: method},
		Suffix:      suff,
		Connection:  conn,
	}
	return
}

func getMessageFromCompoundSignal(arg []string) []byte {
	return []byte(strings.Join(arg, " "))
}

//initializing all the service states to passives
func initializeState(configuration *parser.YAMLBenchRoutesType) {
	configuration.Config.UtilsConf.ServicesSignal.Ping = "passive"
	configuration.Config.UtilsConf.ServicesSignal.Jitter = "passive"
	configuration.Config.UtilsConf.ServicesSignal.FloodPing = "passive"
	configuration.Config.UtilsConf.ServicesSignal.ReqResDelayMonitoring = "passive"
	_, e := configuration.Write()
	if e != nil {
		panic(e)
	}

}
