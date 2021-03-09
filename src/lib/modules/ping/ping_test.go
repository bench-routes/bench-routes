package ping

import (
	"fmt"
	"testing"
	"time"

	"github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/bench-routes/bench-routes/src/lib/filters"
	"github.com/bench-routes/bench-routes/src/lib/utils"
	"github.com/bench-routes/bench-routes/src/lib/utils/prom"
	"github.com/bench-routes/bench-routes/tsdb"
)

var (
	worker               *Ping
	configuration        *parser.Config
	configurationPath    = "../testfiles/configuration.yaml"
	targets              = make(map[string]*utils.MachineType)
	chainSet             = tsdb.NewChainSet(tsdb.FlushAsTime, time.Second*1)
	targetMachineMetrics = prom.MachineMetrics()
)

func initVars() {
	configuration = parser.New(configurationPath)
	configuration.Load().Validate()
	worker = New(configuration, TestInterval{OfType: "sec", Duration: 0}, &targets)
	chainSet.Run()
	for _, r := range configuration.Config.Routes {
		hash := utils.GetHash(filters.HTTPPingFilterValue(r.URL))
		if _, ok := targets[hash]; !ok {
			path := fmt.Sprintf("../testfiles/%s_ping.json", hash)
			targets[hash] = &utils.MachineType{
				IPDomain: filters.HTTPPingFilterValue(r.URL),
				Ping:     tsdb.NewChain(path).Init(),
				Metrics:  targetMachineMetrics,
			}
			chainSet.Register(hash, targets[hash].Ping)
		}
	}
}

func Test_module_PING(T *testing.T) {
	initVars()
	go worker.Iterate("start", true)
	time.Sleep(time.Second * 5)
	go worker.Iterate("stop", true)
}
