package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/zairza-cetb/bench-routes/src/collector/process"
	"github.com/zairza-cetb/bench-routes/src/lib/api"
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/parser"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

var (
	port                    = ":9090" // default listen and serve at 9090
	enableProcessCollection = true    // default collection of process metrices in host of bench-routes
	upgrader                = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}
	configuration parser.YAMLBenchRoutesType
)

func init() {
	logger.Terminal("initializing bench-routes ...", "p")

	configuration.Address = utils.ConfigurationFilePath
	configuration = *configuration.Load()
	configuration.Validate()

	var ConfigURLs []string
	initializeState(&configuration)

	// Build TSDB chain
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
	msg := "initial chain formation time: " + time.Since(p).String()
	logger.Terminal(msg, "p")

	// keep the below line to the end of file so that we ensure that we give a confirmation message only when all the
	// required resources for the application is up and healthy
	logger.Terminal("Bench-routes is up and running", "p")
}

func main() {

	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	api := api.New()
	router := mux.NewRouter()

	router.HandleFunc("/", api.Home)
	router.HandleFunc("/test", api.TestTemplate)
	router.HandleFunc("/service-state", api.ServiceState)
	router.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			msg := "error using upgrader" + err.Error()
			logger.Terminal(msg, "f")
		}

		// capture client request for enabling series of responses unless its killed
		for {
			messageType, message, err := ws.ReadMessage()
			if err != nil {
				logger.Terminal("connection to the terminal lost.", "p")
				logger.Terminal(err.Error(), "p")
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
			msg := "type: " + strconv.Itoa(messageType) + " \n message: " + sig
			logger.Terminal(msg, "p")
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
				querier(ws, inStream, qReqResDelayRoute{})
			}
		}
	})

	if len(os.Args) > 2 && os.Args[2] != "" {
		enableProcessCollection, _ = strconv.ParseBool(os.Args[2])
	}
	if enableProcessCollection {
		go func() {
			const (
				path           = "collector-store/"
				scrapeDuration = time.Second * 15 // default scrape duration for process metrics.
				// TODO: accept scrape-duration for process metrices via args.
			)

			var (
				wg              sync.WaitGroup
				buffer          = process.NewProcessReader()
				collectionCount = 0
			)

			assignChaintoMap := func(c *map[string]*tsdb.Chain, n, path string) {
				(*c)[n] = tsdb.NewChain(path)
				(*c)[n].Init().Commit()
			}

			processChains := make(map[string]*tsdb.Chain)

			for {
				collectionCount++

				if collectionCount%10 != 1 { // in every blocks of 10.
					logger.File(fmt.Sprintf("collection-count: %d; scrape-duration: %fsecs", collectionCount, scrapeDuration.Seconds()), "p")
				} else {
					logger.Terminal(fmt.Sprintf("collection-count: %d; scrape-duration: %fsecs", collectionCount, scrapeDuration.Seconds()), "p")
				}

				if _, err := buffer.UpdateCurrentProcesses(); err != nil {
					panic(err)
				}

				wg.Add(buffer.TotalRunningProcesses)

				for _, ps := range *buffer.ProcessesDetails {
					if processChains[ps.FilteredCommand] == nil {
						p := fmt.Sprintf("%s%s.json", path, ps.FilteredCommand)
						assignChaintoMap(&processChains, ps.FilteredCommand, p)
					}
					b := *tsdb.GetNewBlock("ps", ps.Encode())
					processChains[ps.FilteredCommand].Append(b).Commit()
					wg.Done()
				}

				runtime.GC()

				wg.Wait()
				time.Sleep(scrapeDuration)
			}
		}()
	}

	logger.Terminal(http.ListenAndServe(port, cors.Default().Handler(router)).Error(), "f")

}

func cleanup() {
	logger.Terminal(fmt.Sprintf("Alive %d goroutines", runtime.NumGoroutine()), "p")
	configuration := configuration.Refresh()
	values := reflect.ValueOf(configuration.Config.UtilsConf.ServicesSignal)
	typeOfServiceState := values.Type()
	type serviceState struct {
		service string
		state   string
	}
	serviceStateValues := []serviceState{}

	for i := 0; i < values.NumField(); i++ {
		n := serviceState{service: typeOfServiceState.Field(i).Name, state: values.Field(i).Interface().(string)}
		serviceStateValues = append(serviceStateValues, n)
	}
	for _, node := range serviceStateValues {
		if node.state == "active" {
			switch node.service {
			case "Ping":
				HandlerPingGeneral("stop")
			case "FloodPing":
				HandlerFloodPingGeneral("stop")
			case "Jitter":
				HandlerJitterGeneral("stop")
			case "ReqResDelayMonitoring":
				HandleReqResGeneral("stop")
			}
		}
	}
	logger.Terminal(fmt.Sprintf("Alive %d goroutines after cleaning up.", runtime.NumGoroutine()), "p")
}

func chainInitialiser(chain *[]*tsdb.Chain, conf interface{}, basePath, Type string) {
	msg := "forming " + Type + " chain ... "
	logger.Terminal(msg, "p")
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

	logger.Terminal("finished "+Type+" chain", "p")
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

// initializeState initializes all state values to passive.
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
