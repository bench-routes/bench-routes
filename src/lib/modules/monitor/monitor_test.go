package monitor

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/config"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/src/lib/utils/prom"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

var (
	worker            *Monitor
	configuration     *parser.Config
	chainSet          = tsdb.NewChainSet(tsdb.FlushAsTime, time.Second*1)
	configurationPath = "../testfiles/configuration.yaml"
	targets           = make(map[string]*utils.BRMatrix)
	endpointMetrics   = prom.EndpointMetrics()
)

func initVars() {
	configuration = parser.New(configurationPath)
	configuration.Load().Validate()
	worker = New(configuration, TestInterval{OfType: "sec", Duration: 2}, &targets)
	chainSet.Run()
	for _, r := range configuration.Config.Routes {
		hash := URLHash(r)
		if _, ok := targets[hash]; !ok {
			path := fmt.Sprintf("../testfiles/%s_monitor.json", hash)
			targets[hash] = &utils.BRMatrix{FullURL: r.URL,
				Route:        r,
				MonitorChain: tsdb.NewChain(path).Init(),
				Metrics:      endpointMetrics,
			}
			chainSet.Register(hash, targets[hash].MonitorChain)
		}
	}
}

func Test_module_MONITOR(T *testing.T) {
	initVars()
	go worker.Iterate("start", true)
	time.Sleep(time.Second * 5)
	go worker.Iterate("stop", true)
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
