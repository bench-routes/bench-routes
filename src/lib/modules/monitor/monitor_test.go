package monitor

import (
	"testing"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/config"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

var (
	configuration     *parser.Config
	chain             []*tsdb.Chain
	configurationPath = "../testfiles/configuration.yaml"
)

func initVars() {
	configuration = parser.New(configurationPath)
	configuration.Load().Validate()

	setDefaultServicesState(configuration)

	chain = []*tsdb.Chain{
		tsdb.NewChain("../testfiles/Test_Monitor_1.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_2.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_3.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_4.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_5.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_6.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_7.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_8.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_9.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_10.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_11.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_12.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_13.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_14.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_15.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_16.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_17.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_18.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_19.json").Init(),
		tsdb.NewChain("../testfiles/Test_Monitor_20.json").Init(),
	}
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

func Test_module_MONITOR(T *testing.T) {
	initVars()
	ping := New(configuration, TestInterval{OfType: "min", Duration: 0}, &chain)
	go ping.Iterate("start", true)
	time.Sleep(time.Second * 30)
	go ping.Iterate("stop", true)
}
