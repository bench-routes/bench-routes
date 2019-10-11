package response

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"github.com/zairza-cetb/bench-routes/src/lib/utils/parser"
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
func HandleResponseDelayForRoute(responseChains [][]*tsdb.Chain, route parser.Routes, tsdbNameHash string, wg *sync.WaitGroup) {
	routeSuffix := filters.RouteDestroyer(route.URL)
	// Init paths for request-response-monitoring
	pathDelay := PathReqResDelay + "/" + "chunk_req_res_" + routeSuffix + "_delay.json"
	pathLength := PathReqResDelay + "/" + "chunk_req_res_" + routeSuffix + "_length.json"
	pathStatusCode := PathReqResDelay + "/" + "chunk_req_res_" + routeSuffix + "_status.json"
	c := make(chan utils.Response)
	go RouteDispatcher(route, c)
	responseObject := <-c
	// Store the respective attributes of the
	// response object in the respective TSDB
	log.Printf("Writing responseObject to TSDB for %s", route.URL)

	// Create response delay block to be
	// saved inside TSDB
	blockDelay := tsdb.Block{
		NextBlock:      nil,
		PrevBlock:      nil,
		Datapoint:      float32(responseObject.Delay),
		Timestamp:      time.Now(),
		NormalizedTime: 0,
	}
	for index := range responseChains[0] {
		if responseChains[0][index].Path == pathDelay {
			responseChains[0][index] = responseChains[0][index].Append(blockDelay)
			responseChains[0][index].Save()
			break
		}
	}

	// Create a response length to be
	// saved inside TSDB
	blockLength := tsdb.Block{
		NextBlock:      nil,
		PrevBlock:      nil,
		Datapoint:      float32(responseObject.ResLength),
		Timestamp:      time.Now(),
		NormalizedTime: 0,
	}
	for index := range responseChains[1] {
		if responseChains[1][index].Path == pathLength {
			responseChains[1][index] = responseChains[1][index].Append(blockLength)
			responseChains[1][index].Save()
			break
		}
	}

	// Create a response statusCode to be
	// saved inside TSDB
	blockStatusCode := tsdb.Block{
		NextBlock:      nil,
		PrevBlock:      nil,
		Datapoint:      float32(responseObject.ResStatusCode),
		Timestamp:      time.Now(),
		NormalizedTime: 0,
	}
	for index := range responseChains[2] {
		if responseChains[2][index].Path == pathStatusCode {
			responseChains[2][index] = responseChains[2][index].Append(blockStatusCode)
			responseChains[2][index].Save()
			break
		}
	}

	wg.Done()
}

// RouteDispatcher dispatches a route to respective handlers based on it's request type
func RouteDispatcher(route parser.Routes, c chan utils.Response) {
	if route.Method == "GET" {
		res := HandleGetRequest(route.URL)
		c <- res
	} else {
		// Send a very large integer to automatically rule out as it
		// is much much larger than the threshold
		c <- utils.Response{Delay: math.MaxInt32, ResLength: 0, ResStatusCode: 100}
	}
}

// HandleGetRequest specifically handles routes with GET Requests. Calculates timestamp before
// and after processing of each request and returns the difference
func HandleGetRequest(url string) utils.Response {
	// Time init
	start := time.Now().UnixNano()
	resp := utils.SendGETRequest(url)
	resLength := resp.ContentLength
	respStatusCode := resp.StatusCode
	defer resp.Body.Close()

	end := time.Now().UnixNano()
	diff := int((end - start) / int64(time.Millisecond))

	return utils.Response{Delay: diff, ResLength: resLength, ResStatusCode: respStatusCode}
}
