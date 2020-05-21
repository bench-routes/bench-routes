package main

import (
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
	"github.com/zairza-cetb/bench-routes/src/lib/api"
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/jitter"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/monitor"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/ping"
	"github.com/zairza-cetb/bench-routes/src/lib/parser"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/src/metrics/process"
	sysMetrics "github.com/zairza-cetb/bench-routes/src/metrics/system"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

var (
	port                        = ":9090" // default listen and serve at 9090
	enableProcessCollection     = false   // default collection of process metrics in host of bench-routes
	processCollectionScrapeTime = time.Second * 5
	systemCollectionScrapeTime  = time.Second * 10
	upgrader                    = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}
	conf *parser.YAMLBenchRoutesType
)

const (
	uiPathV11 = "ui-builds/v1.1/"
)

func main() {
	if len(os.Args) > 2 && os.Args[2] != "" {
		enableProcessCollection, _ = strconv.ParseBool(os.Args[2])
		port = ":" + os.Args[1]
	} else if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}

	conf = parser.New(utils.ConfigurationFilePath)
	conf.Load().Validate()
	intervals := conf.Config.Interval

	logger.Terminal("initializing...", "p")
	var ConfigURLs []string
	matrix := make(utils.BRmap)
	setDefaultServicesState(conf)

	setMatrixKey := func(matrix utils.BRmap, url string) utils.BRmap {
		var count int
		for k := range matrix {
			if k == url {
				count++
			}
		}
		if count != 0 {
			if count == 1 {
				delete(matrix, url)
				tmp := fmt.Sprintf("%s-%d", url, count)
				matrix[tmp] = &utils.BRMatrix{
					Domain: tmp,
				}
			}
			tmp := fmt.Sprintf("%s-%d", url, count+1)
			matrix[tmp] = &utils.BRMatrix{
				Domain: tmp,
			}
		} else {
			matrix[url] = &utils.BRMatrix{
				Domain: url,
			}
		}
		return matrix
	}

	for _, r := range conf.Config.Routes {
		var found bool
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
		matrix = setMatrixKey(matrix, r.URL)
	}
	var wg sync.WaitGroup
	p := time.Now()
	wg.Add(4)

	chainSet := tsdb.NewChainSet(tsdb.FlushAsTime, time.Second*50)

	go initialise(&wg, &matrix, chainSet, &utils.Pingc, ConfigURLs, utils.PathPing, "ping")
	go initialise(&wg, &matrix, chainSet, &utils.FPingc, ConfigURLs, utils.PathFloodPing, "flood_ping")
	go initialise(&wg, &matrix, chainSet, &utils.Jitterc, ConfigURLs, utils.PathJitter, "jitter")
	go initialise(&wg, &matrix, chainSet, &utils.RespMonitoringc, conf.Config.Routes, utils.PathReqResDelayMonitoring, "req_res")
	wg.Wait()
	msg := "initialization time: " + time.Since(p).String()
	logger.Terminal(msg, "p")

	chainSet.Run()

	service := struct {
		Ping    *ping.Ping
		Jitter  *jitter.Jitter
		PingF   *ping.FloodPing
		Monitor *monitor.Monitor
	}{
		Ping:    ping.New(conf, ping.TestInterval{OfType: intervals[0].Type, Duration: *intervals[0].Duration}, utils.Pingc),
		Jitter:  jitter.New(conf, jitter.TestInterval{OfType: intervals[0].Type, Duration: *intervals[1].Duration}, utils.Jitterc),
		PingF:   ping.Newf(conf, ping.TestInterval{OfType: intervals[0].Type, Duration: *intervals[0].Duration}, utils.FPingc, conf.Config.Password),
		Monitor: monitor.New(conf, monitor.TestInterval{OfType: intervals[2].Type, Duration: *intervals[2].Duration}, utils.RespMonitoringc),
	}

	api := api.New(&matrix)
	router := mux.NewRouter()

	// Persistent connection for real-time updates between UI and the service.
	router.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			msg := "error using upgrader" + err.Error()
			logger.Terminal(msg, "f")
			os.Exit(1)
		}

		format := func(b bool) []byte {
			return []byte(strconv.FormatBool(b))
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
			logger.File(msg, "p")
			// generate appropriate signals from incoming messages
			switch sig {
			// ping
			case "force-start-ping":
				if e := ws.WriteMessage(1, format(service.Ping.Iterate("start", false))); e != nil {
					panic(e)
				}
			case "force-stop-ping":
				if e := ws.WriteMessage(1, format(service.Ping.Iterate("stop", false))); e != nil {
					panic(e)
				}

				// flood-ping
			case "force-start-flood-ping":
				if e := ws.WriteMessage(1, format(service.PingF.Iteratef("start", false))); e != nil {
					panic(e)
				}
			case "force-stop-flood-ping":
				if e := ws.WriteMessage(1, format(service.PingF.Iteratef("stop", false))); e != nil {
					panic(e)
				}

				// jitter
			case "force-start-jitter":
				if e := ws.WriteMessage(1, format(service.Jitter.Iterate("start", false))); e != nil {
					panic(e)
				}
			case "force-stop-jitter":
				if e := ws.WriteMessage(1, format(service.Jitter.Iterate("start", false))); e != nil {
					panic(e)
				}

				// request-monitor-monitoring
			case "force-start-req-res-monitoring":
				if e := ws.WriteMessage(1, format(service.Monitor.Iterate("start", false))); e != nil {
					panic(e)
				}
			case "force-stop-req-res-monitoring":
				if e := ws.WriteMessage(1, format(service.Monitor.Iterate("stop", false))); e != nil {
					panic(e)
				}

				// Get config routes details
			case "route-details":
				m := conf.Config.Routes
				if e := ws.WriteMessage(1, filters.RouteYAMLtoJSONParser(m)); e != nil {
					panic(e)
				}
			}
		}
	})

	go func() {
		var (
			metrics = sysMetrics.New()
		)

		type metric struct {
			cpu    string
			memory sysMetrics.MemoryStats
			disk   sysMetrics.DiskStats
		}

		chain := tsdb.NewChain("storage/system.json")
		chain.Init()
		chainSet.Register(chain.Name, chain)

		for {
			// collections for cpu, memory and disk run independently and are
			// time dependent. Hence, running these serially will take more
			// time than the actual `systemCollectionScrapeTime`. Hence, the
			// best way is to run them in parallel and get data via channels,
			// such that systemCollectionScrapeTime >= duration(cpu|memory|disk)
			// will meet excepted systemCollectionScrapeTime. Anything other
			// than this will be inaccurate.
			cpu := make(chan string)
			memory := make(chan sysMetrics.MemoryStats)
			disk := make(chan sysMetrics.DiskStats)

			go metrics.GetTotalCPUUsage(cpu)
			go metrics.GetVirtualMemoryStats(memory)
			go metrics.GetDiskIOStats(disk)

			data := &metric{
				cpu:    <-cpu,
				memory: <-memory,
				disk:   <-disk,
			}

			encoded := metrics.Combine(
				metrics.Encode(data.cpu), metrics.Encode(data.memory), metrics.Encode(data.disk),
			)

			block := tsdb.GetNewBlock("sys", encoded)
			chain.Append(*block)
			time.Sleep(systemCollectionScrapeTime)
		}
	}()

	if enableProcessCollection {
		go func() {
			var (
				path           = "collector-store/"
				scrapeDuration = processCollectionScrapeTime // default scrape duration for process metrics.
				// TODO: accept scrape-duration for process metrics via args.

				wg              sync.WaitGroup
				buffer          = process.New()
				collectionCount = 0
			)
			assignChaintoMap := func(c *map[string]*tsdb.Chain, n, path string) {
				(*c)[n] = tsdb.NewChain(path)
				(*c)[n].Init()
				chainSet.Register((*c)[n].Name, (*c)[n])
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
					logger.File(fmt.Sprintf("Fatal: %s", err.Error()), "f")
					os.Exit(1)
				}
				wg.Add(buffer.TotalRunningProcesses)
				for _, ps := range *buffer.ProcessesDetails {
					if processChains[ps.FilteredCommand] == nil {
						p := fmt.Sprintf("%s%s.json", path, ps.FilteredCommand)
						assignChaintoMap(&processChains, ps.FilteredCommand, p)
					}
					b := tsdb.GetNewBlock("ps", ps.Encode())
					processChains[ps.FilteredCommand].Append(*b)
					wg.Done()
				}
				wg.Wait()
				time.Sleep(scrapeDuration)
			}
		}()
	}

	// clean tsdb blocks in regular intervals.
	go func() {
		runtime.GC()
		time.Sleep(time.Duration(time.Minute * 3))
	}()

	// Reset Services.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Terminal(fmt.Sprintf("Alive %d goroutines", runtime.NumGoroutine()), "p")
		conf.Refresh()
		values := reflect.ValueOf(conf.Config.UtilsConf.ServicesSignal)
		typeOfServiceState := values.Type()

		type serviceState struct {
			service string
			state   string
		}

		var serviceStateValues []serviceState
		for i := 0; i < values.NumField(); i++ {
			n := serviceState{service: typeOfServiceState.Field(i).Name, state: values.Field(i).Interface().(string)}
			serviceStateValues = append(serviceStateValues, n)
		}
		for _, node := range serviceStateValues {
			if node.state == "active" {
				switch node.service {
				case "Ping":
					service.Ping.Iterate("stop", false)
				case "FloodPing":
					service.PingF.Iteratef("stop", false)
				case "Jitter":
					service.Jitter.Iterate("stop", false)
				case "ReqResDelayMonitoring":
					service.Monitor.Iterate("stop", false)
				}
			}
		}
		logger.Terminal(fmt.Sprintf("Alive %d goroutines after cleaning up.", runtime.NumGoroutine()), "p")
		os.Exit(0)
	}()

	// API endpoints.
	{
		// static servings.
		{
			router.Handle("/", http.FileServer(http.Dir(uiPathV11)))
			router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(uiPathV11+"assets/"))))
			router.PathPrefix("/manifest.json").Handler(http.StripPrefix("/manifest.json", http.FileServer(http.Dir(uiPathV11+"/manifest.json"))))
			router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(uiPathV11+"static/"))))
		}
		router.HandleFunc("/br-live-check", api.Home)
		router.HandleFunc("/test", api.TestTemplate)
		router.HandleFunc("/service-state", api.ServiceState)
		router.HandleFunc("/routes-summary", api.RoutesSummary)
		router.HandleFunc("/get-route-time-series", api.TSDBPathDetails)
		router.HandleFunc("/query-matrix", api.SendMatrix)
		router.HandleFunc("/query", api.Query)
	}

	logger.Terminal(http.ListenAndServe(port, cors.Default().Handler(router)).Error(), "f")
	// keep the below line to the end of file so that we ensure that we give a confirmation message only when all the
	// required resources for the application is up and healthy.
	logger.Terminal("Bench-routes is up and running", "p")
}

func initialise(wg *sync.WaitGroup, matrix *utils.BRmap, chainSet *tsdb.ChainSet, chain *[]*tsdb.Chain, conf interface{}, basePath, Type string) {
	msg := "forming " + Type + " chain ... "
	logger.File(msg, "p")
	config, ok := conf.([]string)
	if ok {
		for _, v := range config {
			path := basePath + "/chunk_" + Type + "_" + v + ".json"
			resp := tsdb.NewChain(path)
			resp.Init()
			*chain = append(*chain, resp)
			for k := range *matrix {
				index := k
				filters.HTTPPingFilter(&k)
				if k == v {
					switch Type {
					case "ping":
						(*matrix)[index].PingChain = resp
					case "jitter":
						(*matrix)[index].JitterChain = resp
					case "flood_ping":
						(*matrix)[index].FPingChain = resp
					}
					break
				}
			}
		}
	}
	if configRes, ok := conf.([]parser.Routes); ok {
		for _, v := range configRes {
			path := basePath + "/chunk_" + Type + "_" + filters.RouteDestroyer(v.URL+"_"+v.Route) + ".json"
			resp := tsdb.NewChain(path)
			resp.Init()
			*chain = append(*chain, resp)
			for k := range *matrix {
				filters.HTTPPingFilter(&v.URL)
				if k == v.URL {
					(*matrix)[k] = &utils.BRMatrix{
						URL:          v.URL,
						Method:       v.Method,
						Route:        v.Route,
						Headers:      v.Header,
						Params:       v.Params,
						MonitorChain: resp,
						Domain:       fmt.Sprintf("%s: %s/%s", v.Method, v.URL, v.Route),
						PingChain:    (*matrix)[k].PingChain,
						JitterChain:  (*matrix)[k].JitterChain,
						FPingChain:   (*matrix)[k].FPingChain,
					}
					break
				}
			}
		}
	}
	for _, chain := range *chain {
		chainSet.Register(chain.Name, chain)
	}

	logger.Terminal("finished "+Type+" chain", "p")
	wg.Done()
}

// setDefaultServicesState initializes all state values to passive.
func setDefaultServicesState(configuration *parser.YAMLBenchRoutesType) {
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
