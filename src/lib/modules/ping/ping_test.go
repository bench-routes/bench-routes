package ping

import (
	"testing"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/parser"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

var (
	configuration     *parser.YAMLBenchRoutesType
	chainPing         []*tsdb.Chain
	configurationPath = "../testfiles/configuration.yaml"
)

func initVars() {
	configuration = parser.New(configurationPath)
	configuration.Load().Validate()

	setDefaultServicesState(configuration)

	chainPing = []*tsdb.Chain{
		tsdb.NewChain("../testfiles/Test_Ping_1.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_2.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_3.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_4.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_5.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_6.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_7.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_8.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_9.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_10.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_11.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_12.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_13.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_14.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_15.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_16.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_17.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_18.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_19.json").Init().Commit(),
		tsdb.NewChain("../testfiles/Test_Ping_20.json").Init().Commit(),
	}
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

func Test_module_PING(T *testing.T) {
	initVars()
	ping := New(configuration, TestInterval{OfType: "min", Duration: 0}, chainPing)
	go ping.Iterate("start", true)
	time.Sleep(time.Second * 30)
	go ping.Iterate("stop", true)
}

func Test_module_FLOOD_PING(T *testing.T) {
	initVars()
	ping := Newf(configuration, TestInterval{OfType: "min", Duration: 0}, chainPing, "")
	go ping.Iteratef("start", true)
	time.Sleep(time.Second * 30)
	go ping.Iteratef("stop", true)
}
