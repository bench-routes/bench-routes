package lib

import (
	"fmt"
	"github.com/zairza-cetb/bench-routes/src/lib/modules/ping"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"log"
	"sync"
	"time"
)

const (
	// ConfigurationFilePath is the constant path to the configuration file needed to start the application
	// written from root file since the application starts from `make run`
	ConfigurationFilePath = "storage/local-config.yml"
)

var (
	// Configuration corresponds to the global use of configuration settings
	Configuration utils.YAMLBenchRoutesType
)

func init() {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)

	// load configuration file
	Configuration.Address = ConfigurationFilePath
	Configuration = *Configuration.Load()

	// keep the below line to the end of file so that we ensure that we give a confirmation message only when all the
	// required resources for the application is up and healthy
	log.Println("Bench-routes is up and running")
}

// HandlerPingGeneral handles the ping route
func HandlerPingGeneral(signal string) bool {

	// Get latest service state settings
	Configuration = Configuration.Refresh()
	pingServiceState := Configuration.Config.UtilsConf.ServicesSignal.Ping

	switch signal {
	case "start":
		if pingServiceState == "passive" {

			Configuration.Config.UtilsConf.ServicesSignal.Ping = "active"
			_, e := Configuration.Write()
			if e != nil {
				panic(e)
			}
			go func() {
				fmt.Println("inside goroutine")
				pingConfig := Configuration.Config.Routes
				pingInterval := getInterval(Configuration.Config.Interval, "ping")
				if pingInterval == (testInterval{}) {
					log.Fatalf("interval not found in configuration file for ping")
					return
				}
				fmt.Println("ping configuration is")
				fmt.Println(pingConfig)

				var urlStack []string
				for _, route := range pingConfig {
					url := route.URL
					exists := false
					for _, e := range urlStack {
						if e == url {
							exists = true
						}
					}
					if !exists {
						urlStack = append(urlStack, url)
					}
				}
				for {
					fmt.Println("inside for infinite")
					Configuration = Configuration.Refresh()

					switch Configuration.Config.UtilsConf.ServicesSignal.Ping {
					case "active":
						var wg sync.WaitGroup
						wg.Add(len(urlStack))
						for u := range urlStack {
							go ping.HandlePing(&urlStack[u], 10, &wg)
						}
						wg.Wait()
					case "passive":
						// terminate the goroutine
						log.Printf("terminating ping goroutine\n")
						return
					default:
						log.Fatalf("invalid service-state value of ping\n")
						return
					}

					intrv := time.Duration(pingInterval.duration)
					switch pingInterval.ofType {
					case "hr":
						time.Sleep(intrv * time.Hour)
					case "min":
						time.Sleep(intrv * time.Minute)
					case "sec":
						time.Sleep(intrv * time.Second)
					default:
						log.Fatalf("invalid interval-type for ping\n")
						return
					}
				}
			}()
			return true
		}
	case "stop":
		Configuration.Config.UtilsConf.ServicesSignal.Ping = "passive"
		_, e := Configuration.Write()
		if e != nil {
			panic(e)
		}
		return true
	default:
		log.Fatalf("invalid signal")
	}
	return false
}

type testInterval struct {
	ofType   string
	duration int64
}

func getInterval(intervals []utils.Interval, testName string) testInterval {
	fmt.Println("here")
	fmt.Println(intervals)
	for _, intrv := range intervals {
		if testName == intrv.Test {
			return testInterval{
				ofType:   intrv.Type,
				duration: intrv.Duration,
			}
		}
	}

	// if the execution reaches this then it is an error as the interval was not found in the configuration file
	log.Panicf("interval not found in the configuration file\n")
	return testInterval{}
}
