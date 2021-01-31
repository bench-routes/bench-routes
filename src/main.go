package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"github.com/zairza-cetb/bench-routes/tsdb/v1"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/zairza-cetb/bench-routes/src/lib/api"
	parser "github.com/zairza-cetb/bench-routes/src/lib/config"
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/jitter"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/monitor"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/ping"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/src/lib/utils/prom"
	"github.com/zairza-cetb/bench-routes/src/metrics/journal"
	"github.com/zairza-cetb/bench-routes/src/metrics/process"
	sysMetrics "github.com/zairza-cetb/bench-routes/src/metrics/system"
)

var (
	port                        = ":9990" // default listen and serve at 9090
	enableProcessCollection     = false   // default collection of process metrics in host of bench-routes
	processCollectionScrapeTime = time.Second * 5
	defaultScrapeTime           = time.Second * 3
	systemMetricsPath           = "storage/system.json"
	journalMetricsPath          = "storage/journal.json"
	// matrix is a collection (as map) of instances where each
	// instance is composed of ping, jitter, flood-ping and monitor
	// chain paths. matrix is used in the monitoring screen to
	// reduce the http request by grouping them based on routes.
	// Filters out identical routes that arrive to be monitored.
	matrix   = make(map[string]*utils.BRMatrix)
	reload   = make(chan struct{})
	start = make(chan struct{})
	done     = make(chan struct{})
	conf     *parser.Config
	chainSet = v1.NewChainSet(v1.FlushAsTime, time.Second*300)
	// targetMachineCalc contains calculations that are machine/vm/load-balancer
	// specific. These involve use of IP addresses/Domain names respectively.
	// We use this to keep track of host names because all the routes that share
	// the same hostname will have the same ping, jitter and flood-ping.
	targetMachineCalc     = make(map[string]*utils.MachineType)
	targetMachineMetrics  = prom.MachineMetrics()
	endpointMetrics       = prom.EndpointMetrics()
	configurationFilePath = "local-config.yml"
	// Storage paths.
	storageDir     = "storage"
	pathPing       = fmt.Sprintf("%s/ping", storageDir)
	pathJitter     = fmt.Sprintf("%s/jitter", storageDir)
	pathFloodPing  = fmt.Sprintf("%s/flood-ping", storageDir)
	pathMonitoring = fmt.Sprintf("%s/monitoring", storageDir)
)

func main() {
	if len(os.Args) > 2 && os.Args[2] != "" {
		enableProcessCollection, _ = strconv.ParseBool(os.Args[2])
		port = ":" + os.Args[1]
	} else if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}
	log.Infoln("initializing bench routes...")
	conf = parser.New(configurationFilePath)
	conf.Load().Validate()
	setDefaultServicesState(conf)
	intervals := conf.Config.Interval
	// Loads workers on to the memory
	workers := &struct {
		Ping    *ping.Ping
		Jitter  *jitter.Jitter
		PingF   *ping.FloodPing
		Monitor *monitor.Monitor
	}{
		Ping:    ping.New(conf, ping.TestInterval{OfType: intervals[0].Type, Duration: *intervals[0].Duration}, &targetMachineCalc),
		Jitter:  jitter.New(conf, jitter.TestInterval{OfType: intervals[0].Type, Duration: *intervals[1].Duration}, &targetMachineCalc),
		PingF:   ping.Newf(conf, ping.TestInterval{OfType: intervals[0].Type, Duration: *intervals[0].Duration}, conf.Config.Password, &targetMachineCalc),
		Monitor: monitor.New(conf, monitor.TestInterval{OfType: intervals[2].Type, Duration: *intervals[2].Duration}, &matrix),
	}
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	prometheus.MustRegister(
		targetMachineMetrics.Ping,
		targetMachineMetrics.PingCount,
		targetMachineMetrics.Jitter,
		targetMachineMetrics.JitterCount,
		targetMachineMetrics.FPing,
		targetMachineMetrics.FPingCount,
		endpointMetrics.ResponseDelay,
		endpointMetrics.ResponseLength,
		endpointMetrics.StatusCode,
		endpointMetrics.MonitorCount,
	)
	// go routine that reloads the entire map(cMap).
	go func() {
		fmt.Println("Start goroutine...")
		select {
			case route := <- reload:
				fmt.Println("Reload goroutine...")
				// add the exact route to
				// be monitored to the chainSet
				fmt.Println(route)
			case <- start:
				fmt.Println("Start bench routes -> -> ->")
				conf.Refresh()
				p := time.Now()
				// Iterate all the routes present in
				// the config file.
				for _, r := range conf.Config.Routes {
					hash := URLHash(r)
					// If the urlHash is not present in
					// matrix, then most likely it has been
					// added recently, we just create a new record
					// for it in the matrix.
					if _, ok := matrix[hash]; !ok {
						var (
							pathPing      = fmt.Sprintf("%s/chunk_ping_%s.json", pathPing, hash)
							pathJitter    = fmt.Sprintf("%s/chunk_jitter_%s.json", pathJitter, hash)
							pathFloodPing = fmt.Sprintf("%s/chunk_flood_ping_%s.json", pathFloodPing, hash)
							pathMonitor   = fmt.Sprintf("%s/chunk_monitor_%s.json", pathMonitoring, hash)
						)
						uHash := utils.GetHash(filters.HTTPPingFilterValue(r.URL))
						if _, ok := targetMachineCalc[uHash]; !ok {
							targetMachineCalc[uHash] = &utils.MachineType{
								IPDomain: filters.HTTPPingFilterValue(r.URL),
								Ping:     v1.NewChain(pathPing).Init(),
								Jitter:   v1.NewChain(pathJitter).Init(),
								FPing:    v1.NewChain(pathFloodPing).Init(),
								Metrics:  targetMachineMetrics,
							}
							// commits new chain after every interval.
							chainSet.Register(fmt.Sprintf("%s-ping", uHash), targetMachineCalc[uHash].Ping)
							chainSet.Register(fmt.Sprintf("%s-jitter", uHash), targetMachineCalc[uHash].Jitter)
							chainSet.Register(fmt.Sprintf("%s-fping", uHash), targetMachineCalc[uHash].FPing)
						}
						// finally, adds the route to the matrix
						matrix[hash] = &utils.BRMatrix{
							FullURL:      r.URL,
							Route:        r,
							PingChain:    targetMachineCalc[uHash].Ping,
							JitterChain:  targetMachineCalc[uHash].Jitter,
							FPingChain:   targetMachineCalc[uHash].FPing,
							MonitorChain: v1.NewChain(pathMonitor).Init(),
							Metrics:      endpointMetrics,
						}
						chainSet.Register(fmt.Sprintf("%s-monitor", uHash), matrix[hash].MonitorChain)
					}
					// If the route is already present in the matrix,
					// the route is already monitored, we just continue
					// normal execution.
				}
				log.Infoln("initialization time: " + time.Since(p).String())
				done <- struct{}{}
			default:
			fmt.Println("Default goroutine...")
		}
	}()
	start <- struct{}{}
	<-done
	chainSet.Run()

	apiInstance := api.New(&matrix, conf, configurationFilePath, workers, reload, done)
	router := mux.NewRouter()
	apiInstance.Register(router)

	go func() {
		metrics := sysMetrics.New()
		type metric struct {
			cpu    *string
			memory *sysMetrics.MemoryStats
			disk   *sysMetrics.DiskStats
			net    *sysMetrics.NetworkStats
		}

		chain := v1.NewChain(systemMetricsPath)
		p := time.Now()
		chain.Init()
		log.Infoln("initialized system-metrics...", time.Since(p))
		chainSet.Register(chain.Name, chain)

		for {
			// collections for cpu, memory and disk run independently and are
			// time dependent. Hence, running these serially will take more
			// time than the actual `defaultScrapeTime`. Hence, the
			// best way is to run them in parallel and get data via channels,
			// such that defaultScrapeTime >= duration(cpu|memory|disk)
			// will meet excepted defaultScrapeTime. Anything other
			// than this will be inaccurate.
			cpu := make(chan *string)
			memory := make(chan *sysMetrics.MemoryStats)
			disk := make(chan *sysMetrics.DiskStats)
			net := make(chan *sysMetrics.NetworkStats)

			go metrics.GetTotalCPUUsage(cpu)
			go metrics.GetVirtualMemoryStats(memory)
			go metrics.GetDiskIOStats(disk)
			go metrics.GetNetworkStats(net)

			data := &metric{
				cpu:    <-cpu,
				memory: <-memory,
				disk:   <-disk,
				net:    <-net,
			}
			encoded := metrics.Combine(
				metrics.Encode(*data.cpu), metrics.Encode(*data.memory), metrics.Encode(*data.disk), metrics.Encode(*data.net),
			)

			block := v1.GetNewBlock("sys", encoded)
			chain.Append(*block)
			time.Sleep(defaultScrapeTime)
		}
	}()

	if !(runtime.GOOS == "windows" || runtime.GOOS == "darwin") {
		go func() {
			metrics := journal.New()
			chain := v1.NewChain(journalMetricsPath)
			p := time.Now()
			chain.Init()
			log.Infoln("initialized journal-metrics...", time.Since(p))
			chainSet.Register(chain.Name, chain)

			for {
				data := metrics.Run().Get()
				datapoint := data.Encode()
				block := v1.GetNewBlock("journal", *datapoint)
				chain.Append(*block)
				time.Sleep(defaultScrapeTime)
			}
		}()
	}

	if enableProcessCollection {
		go func() {
			var (
				path           = "collector-store/"
				scrapeDuration = processCollectionScrapeTime // default scrape duration for process metrics.
				// TODO: accept scrape-duration for process metrics via args.

				wg     sync.WaitGroup
				buffer = process.New()
			)
			assignChaintoMap := func(c *map[string]*v1.Chain, n, path string) {
				(*c)[n] = v1.NewChain(path)
				(*c)[n].Init()
				chainSet.Register((*c)[n].Name, (*c)[n])
			}
			processChains := make(map[string]*v1.Chain)
			for {
				if _, err := buffer.UpdateCurrentProcesses(); err != nil {
					log.Infoln(fmt.Sprintf("Fatal: %s", err.Error()), "f")
					os.Exit(1)
				}
				wg.Add(buffer.TotalRunningProcesses)
				for _, ps := range *buffer.ProcessesDetails {
					if processChains[ps.FilteredCommand] == nil {
						p := fmt.Sprintf("%s%s.json", path, ps.FilteredCommand)
						assignChaintoMap(&processChains, ps.FilteredCommand, p)
					}
					b := v1.GetNewBlock("ps", ps.Encode())
					processChains[ps.FilteredCommand].Append(*b)
					wg.Done()
				}
				wg.Wait()
				time.Sleep(scrapeDuration)
			}
		}()
	}

	// Reset Services.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Infof("Alive %d goroutines", runtime.NumGoroutine())
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
					workers.Ping.Iterate("stop", false)
				case "FloodPing":
					workers.PingF.Iteratef("stop", false)
				case "Jitter":
					workers.Jitter.Iterate("stop", false)
				case "ReqResDelayMonitoring":
					workers.Monitor.Iterate("stop", false)
				}
			}
		}
		log.Infof("Alive %d goroutines after cleaning up.\n", runtime.NumGoroutine())
		os.Exit(0)
	}()
	log.Infoln("Bench-routes is up and running")
	log.Infoln(http.ListenAndServe(port, cors.Default().Handler(router)).Error())
	// keep the below line to the end of file so that we ensure that we give a confirmation message only when all the
	// required resources for the application is up and healthy.
	log.Infoln("Bench-routes is up and running")
}

// setDefaultServicesState initializes all state values to passive.
func setDefaultServicesState(configuration *parser.Config) {
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

// URLHash hashes the passed route to a unique value.
func URLHash(route parser.Route) string {
	var (
		method    = route.Method
		URL       = route.URL
		body      = route.Body
		headers   = route.Header
		params    = route.Params
		hashInput = fmt.Sprintf("%s%s", method, URL)
	)
	mBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	mHeaders, err := json.Marshal(headers)
	if err != nil {
		panic(err)
	}
	mParams, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	hashInput += fmt.Sprintf("%s%s%s", mBody, mHeaders, mParams)
	return utils.GetHash(hashInput)
}
