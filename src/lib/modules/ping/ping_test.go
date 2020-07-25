package ping

import (
	"fmt"
	"testing"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/config"
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

var (
	worker            *Ping
	configuration     *parser.Config
	configurationPath = "../testfiles/configuration.yaml"
	targets           = make(map[string]*utils.MachineType)
	chainSet          = tsdb.NewChainSet(tsdb.FlushAsTime, time.Second*1)
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
			targets[hash] = &utils.MachineType{IPDomain: filters.HTTPPingFilterValue(r.URL), Ping: tsdb.NewChain(path).Init()}
			chainSet.Register(hash, targets[hash].Ping)
		}
	}
}

func Test_module_PING(T *testing.T) {
	initVars()
	go worker.Iterate("start", true)
	time.Sleep(time.Second * 30)
	go worker.Iterate("stop", true)
}
