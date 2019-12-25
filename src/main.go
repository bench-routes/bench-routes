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
	"sync"
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
	testFilesDir  = "tests/"
)

func init() {
	go setupLogger()
	log.Printf("initializing bench-routes ...")

	// load configuration file
	configuration.Address = utils.ConfigurationFilePath
	configuration = *configuration.Load()
	configuration.Validate()

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
	var wg sync.WaitGroup
	p := time.Now()
	wg.Add(4)

	go func() {
		chainInitialiser(&utils.GlobalPingChain, ConfigURLs, utils.PathPing, "ping")
		wg.Done()
	}()

	go func() {
		chainInitialiser(&utils.GlobalFloodPingChain, ConfigURLs, utils.PathFloodPing, "flood_ping")
		wg.Done()
	}()

	go func() {
		chainInitialiser(&utils.GlobalChain, ConfigURLs, utils.PathJitter, "jitter")
		wg.Done()
	}()

	go func() {
		chainInitialiser(&utils.GlobalReqResDelChain, configuration.Config.Routes, utils.PathReqResDelayMonitoring, "req_res")
		wg.Done()
	}()

	wg.Wait()
	log.Printf("initial chain formation time: %s\n", time.Since(p).String())

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
				querier(ws, inStream, qPingRoute{})

			case "Qjitter-route":
				querier(ws, inStream, qJitterRoute{})

			case "Qflood-ping-route":
				querier(ws, inStream, qFloodPingRoute{})

			// TODO just like flood-ping, jitter, ping
			// Querrier signal for Request-response delay
			case "Qrequest-response-delay-route":
				// compMessage := getMessageFromCompoundSignal(inStream[1:])
				// inst := qReqResDelayRoute{}
				// if e := json.Unmarshal(compMessage, &inst); e != nil {
				// 	panic(e)
				// }

				// url := inst.URL
				// method := inst.Method
				// // Gets the Querrier for request-response delay
				// // TODO: Send the method along with URL
				// ql := getQuerier(ws, "req-res-delay", url, method, "_delay")
				// go ql.FetchAllSeriesStringified()
				querier(ws, inStream, qReqResDelayRoute{})
			}
		}
	})

	// launch service
	log.Fatal(http.ListenAndServe(port, nil))

}

func chainInitialiser(chain *[]*tsdb.Chain, conf interface{}, basePath, Type string) {
	log.Printf("forming %s chain ... \n", Type)
	config, ok := conf.([]string)
	if ok {
		for _, v := range config {
			path := basePath + "/chunk_" + Type + "_" + v + ".json"

			resp := &tsdb.Chain{
				Path:           path,
				Chain:          []tsdb.Block{},
				LengthElements: 0,
				Size:           0,
			}
			resp.Init().Commit()
			*chain = append(*chain, resp)
		}
	}
	configRes, ok := conf.([]parser.Routes)
	if ok {
		for _, v := range configRes {
			fmt.Println(v.URL + v.Route)
			path := basePath + "/chunk_" + Type + "_" + filters.RouteDestroyer(v.URL+"_"+v.Route) + ".json"

			resp := &tsdb.Chain{
				Path:           path,
				Chain:          []tsdb.Block{},
				LengthElements: 0,
				Size:           0,
			}
			resp.Init().Commit()
			*chain = append(*chain, resp)
		}
	}

	log.Printf("finished %s chain\n", Type)
}

func querier(ws *websocket.Conn, inComingStream []string, route interface{}) {
	message := getMessageFromCompoundSignal(inComingStream[1:])
	var response []interface{}
	switch route.(type) {
	case qPingRoute:
		inst := qPingRoute{}
		if e := json.Unmarshal(message, &inst); e != nil {
			panic(e)
		}

		raw := getInBlocks(ws, "ping", inst.URL)
		for i, b := range raw {
			decRaw := utils.Decode(b)
			dec, ok := decRaw.(utils.Ping)
			if !ok {
				panic("invalid interface type")
			}
			response = append(response, utils.PingResp{
				Min:            dec.Min,
				Mean:           dec.Mean,
				Max:            dec.Max,
				MDev:           dec.MDev,
				NormalizedTime: b.GetNormalizedTime(),
				Timestamp:      b.GetTimeStamp(),
				Relative:       i,
			})
		}

	case qJitterRoute:
		inst := qJitterRoute{}
		if e := json.Unmarshal(message, &inst); e != nil {
			panic(e)
		}

		raw := getInBlocks(ws, "jitter", inst.URL)
		for i, b := range raw {
			decRaw, ok := utils.Decode(b).(float64)
			if !ok {
				panic("invalid interface type")
			}
			response = append(response, utils.JitterResp{
				Datapoint:      decRaw,
				NormalizedTime: b.GetNormalizedTime(),
				Timestamp:      b.GetTimeStamp(),
				Relative:       i,
			})
		}

	case qFloodPingRoute:
		inst := qFloodPingRoute{}
		if e := json.Unmarshal(message, &inst); e != nil {
			panic(e)
		}

		raw := getInBlocks(ws, "flood-ping", inst.URL)
		for i, b := range raw {
			dec, ok := utils.Decode(b).(utils.FloodPing)
			if !ok {
				panic("invalid interface type")
			}
			response = append(response, utils.FloodPingResp{
				Min:            dec.Min,
				Mean:           dec.Mean,
				Max:            dec.Max,
				MDev:           dec.MDev,
				PacketLoss:     dec.PacketLoss,
				NormalizedTime: b.GetNormalizedTime(),
				Timestamp:      b.GetTimeStamp(),
				Relative:       i,
			})
		}

	case qReqResDelayRoute:
		inst := qReqResDelayRoute{}
		if e := json.Unmarshal(message, &inst); e != nil {
			panic(e)
		}

		raw := getInBlocks(ws, "req-res-delay", inst.URL)
		for i, b := range raw {
			dec, ok := utils.Decode(b).(utils.Response)
			if !ok {
				panic("invalid interface type")
			}
			response = append(response, utils.ResponseResp{
				ResLength:      dec.ResLength,
				ResStatusCode:  dec.ResStatusCode,
				Delay:          dec.Delay,
				NormalizedTime: b.GetNormalizedTime(),
				Timestamp:      b.GetTimeStamp(),
				Relative:       i,
			})
		}
	}
	respond(ws, response)
}

func getInBlocks(ws *websocket.Conn, Type, URL string) []tsdb.Block {
	ql := getQuerier(ws, Type, URL, "", "")
	return inBlocks(ql.FetchAllSeriesStringified())
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
	configuration.Config.UtilsConf.ServicesSignal = parser.ServiceSignals{
		Ping:                  "passive",
		Jitter:                "passive",
		FloodPing:             "passive",
		ReqResDelayMonitoring: "passive",
	}
	if _, e := configuration.Write(); e != nil {
		panic(e)
	}
}

func inBlocks(s string) (tmp []tsdb.Block) {
	if err := json.Unmarshal([]byte(s), &tmp); err != nil {
		panic(err)
	}
	return
}

func respond(ws *websocket.Conn, inf interface{}) {
	js, err := json.Marshal(inf)
	if err != nil {
		panic(err)
	}
	if e := ws.WriteMessage(1, js); e != nil {
		panic(e)
	}
}
