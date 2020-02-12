package response

import (
	"io/ioutil"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/parser"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/tsdb"
)

// Path to storage in tsdb for Req res delay
const (
	// PathPing stores the defualt address of storage directory of ping data
	PathReqResDelay = "storage/req-res-delay-monitoring"
)

// HandleResponseDelayForRoute is the initial entrypoint function for this module which takes
// in a Route struct and supplies it to a function in turn to handle it accordingly. We create
// channels to run tests for each route in parallel, speeding up the process
func HandleResponseDelayForRoute(responseChains []*tsdb.Chain, route parser.Routes, tsdbNameHash string, wg *sync.WaitGroup) {
	routeSuffix := filters.RouteDestroyer(route.URL + "_" + route.Route)
	// Init paths for request-response-monitoring
	path := PathReqResDelay + "/" + "chunk_req_res_" + routeSuffix + ".json"
	c := make(chan utils.Response)
	responseObject := RouteDispatcher(route, c)
	g := getNormalizedBlockString(responseObject)
	block := *tsdb.GetNewBlock("req-res", g)
	for index := range responseChains {
		if responseChains[index].Path == path {
			responseChains[index] = responseChains[index].Append(block)
			responseChains[index].Commit()
			break
		}
	}

	wg.Done()
}

// RouteDispatcher dispatches a route to respective handlers based on it's request type
func RouteDispatcher(route parser.Routes, c chan utils.Response) utils.Response {
	if route.Method == "GET" {
		return HandleGetRequest(route.URL)
	}
	// If fail, then
	// send a very large integer to automatically rule out as it
	// is much much larger than the threshold
	return utils.Response{Delay: math.MaxInt32, ResLength: 0, ResStatusCode: 100}
}

// HandleGetRequest specifically handles routes with GET Requests. Calculates timestamp before
// and after processing of each request and returns the difference
func HandleGetRequest(url string) utils.Response {
	// Time init
	start := time.Now().UnixNano()
	resp := *utils.SendGETRequest(url)
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	resLength := len(content)
	respStatusCode := resp.StatusCode

	end := time.Now().UnixNano()
	diff := int((end - start) / int64(time.Millisecond))
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}

	return utils.Response{Delay: diff, ResLength: resLength, ResStatusCode: respStatusCode}
}

// returns the stringified form of the combined data
func getNormalizedBlockString(b utils.Response) string {
	return strconv.Itoa(b.Delay) + tsdb.BlockDataSeparator + strconv.Itoa(b.ResLength) + tsdb.BlockDataSeparator +
		strconv.Itoa(b.ResStatusCode)
}
