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
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/zairza-cetb/bench-routes/src/lib/api"
	parser "github.com/zairza-cetb/bench-routes/src/lib/config"
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/jitter"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/monitor"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/ping"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/src/metrics/journal"
	"github.com/zairza-cetb/bench-routes/src/metrics/process"
	sysMetrics "github.com/zairza-cetb/bench-routes/src/metrics/system"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

var (
	port                        = ":9090" // default listen and serve at 9090
	enableProcessCollection     = false   // default collection of process metrics in host of bench-routes
	processCollectionScrapeTime = time.Second * 5
	defaultScrapeTime           = time.Second * 3
	systemMetricsPath           = "storage/system.json"
	journalMetricsPath          = "storage/journal.json"
	// matrix is a collection (as map) of instances where each
	// instance is composed of ping, jitter, flood-ping and monitor
	// chain paths. matrix is used in the monitoring screen to
	// reduce the http request by grouping them based on routes.
	// Without matrix, the http traffic would increase 4 times
	// the current count.
	matrix   = make(map[string]*utils.BRMatrix)
	reload   = make(chan struct{})
	done     = make(chan struct{})
	conf     *parser.Config
	chainSet = tsdb.NewChainSet(tsdb.FlushAsTime, time.Second*10)
	// targetMachineCalc contains calculations that are machine/vm/load-balancer
	// specific. These involve use of IP addresses/Domain names respectively.
	targetMachineCalc = make(map[string]*utils.MachineType)
)

func main() {
	if len(os.Args) > 2 && os.Args[2] != "" {
		enableProcessCollection, _ = strconv.ParseBool(os.Args[2])
		port = ":" + os.Args[1]
	} else if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}

	logger.Terminal("initializing...", "p")
	conf = parser.New(utils.ConfigurationFilePath)
	conf.Load().Validate()
	setDefaultServicesState(conf)
	intervals := conf.Config.Interval
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
	go func() {
		for {
			<-reload
			fmt.Println("reloading...")
			conf.Refresh()
			p := time.Now()
			for _, r := range conf.Config.Routes {
				hash := URLHash(r)
				if _, ok := matrix[hash]; !ok {
					var (
						pathPing      = fmt.Sprintf("%s/chunk_ping_%s.json", utils.PathPing, hash)
						pathJitter    = fmt.Sprintf("%s/chunk_jitter_%s.json", utils.PathJitter, hash)
						pathFloodPing = fmt.Sprintf("%s/chunk_flood_ping_%s.json", utils.PathFloodPing, hash)
						pathMonitor   = fmt.Sprintf("%s/chunk_monitor_%s.json", utils.PathMonitoring, hash)
					)
					uHash := utils.GetHash(filters.HTTPPingFilterValue(r.URL))
					if _, ok := targetMachineCalc[uHash]; !ok {
						targetMachineCalc[uHash] = &utils.MachineType{
							IPDomain: filters.HTTPPingFilterValue(r.URL),
							Ping:     tsdb.NewChain(pathPing).Init(),
							Jitter:   tsdb.NewChain(pathJitter).Init(),
							FPing:    tsdb.NewChain(pathFloodPing).Init(),
						}
						chainSet.Register(fmt.Sprintf("%s-ping", uHash), targetMachineCalc[uHash].Ping)
						chainSet.Register(fmt.Sprintf("%s-jitter", uHash), targetMachineCalc[uHash].Jitter)
						chainSet.Register(fmt.Sprintf("%s-fping", uHash), targetMachineCalc[uHash].FPing)
					}
					matrix[hash] = &utils.BRMatrix{
						FullURL:      r.URL,
						Route:        r,
						PingChain:    targetMachineCalc[uHash].Ping,
						JitterChain:  targetMachineCalc[uHash].Jitter,
						FPingChain:   targetMachineCalc[uHash].FPing,
						MonitorChain: tsdb.NewChain(pathMonitor).Init(),
					}
					chainSet.Register(fmt.Sprintf("%s-monitor", uHash), matrix[hash].MonitorChain)
				}
			}
			msg := "initialization time: " + time.Since(p).String()
			logger.Terminal(msg, "p")
			done <- struct{}{}
		}
	}()
	reload <- struct{}{}
	<-done
	chainSet.Run()

	apiInstance := api.New(&matrix, conf, workers, reload, done)
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

		chain := tsdb.NewChain(systemMetricsPath)
		p := time.Now()
		chain.Init()
		fmt.Println("initialized system-metrics...", time.Since(p))
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

			block := tsdb.GetNewBlock("sys", encoded)
			chain.Append(*block)
			time.Sleep(defaultScrapeTime)
		}
	}()

	if !(runtime.GOOS == "windows" || runtime.GOOS == "darwin") {
		go func() {
			metrics := journal.New()
			chain := tsdb.NewChain(journalMetricsPath)
			p := time.Now()
			chain.Init()
			fmt.Println("initialized journal-metrics...", time.Since(p))
			chainSet.Register(chain.Name, chain)

			for {
				data := metrics.Run().Get()
				datapoint := data.Encode()
				block := tsdb.GetNewBlock("journal", *datapoint)
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
		time.Sleep(time.Minute * 3)
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
		logger.Terminal(fmt.Sprintf("Alive %d goroutines after cleaning up.", runtime.NumGoroutine()), "p")
		os.Exit(0)
	}()
	logger.Terminal(http.ListenAndServe(port, cors.Default().Handler(router)).Error(), "f")
	// keep the below line to the end of file so that we ensure that we give a confirmation message only when all the
	// required resources for the application is up and healthy.
	logger.Terminal("Bench-routes is up and running", "p")
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
