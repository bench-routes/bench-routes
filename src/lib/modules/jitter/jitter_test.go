package jitter

import (
	"testing"
	"time"

	parser "github.com/zairza-cetb/bench-routes/src/lib/config"
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
		tsdb.NewChain("../testfiles/Test_Jitter_1.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_2.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_3.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_4.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_5.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_6.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_7.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_8.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_9.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_10.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_11.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_12.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_13.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_14.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_15.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_16.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_17.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_18.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_19.json").Init(),
		tsdb.NewChain("../testfiles/Test_Jitter_20.json").Init(),
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

func Test_module_JITTER(T *testing.T) {
	initVars()
	jitter := New(configuration, TestInterval{OfType: "min", Duration: 0}, &chain)
	go jitter.Iterate("start", true)
	time.Sleep(time.Second * 30)
	go jitter.Iterate("stop", true)
}
