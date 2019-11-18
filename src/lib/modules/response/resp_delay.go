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

// Response struct
// This is the object that we return from resp_delay module
// Contains delay in response and the response length
type Response struct {
	Delay         int
	ResLength     int64
	ResStatusCode int
}

// HandleResponseDelayForRoute is the initial entrypoint function for this module which takes
// in a Route struct and supplies it to a function in turn to handle it accordingly. We create
// channels to run tests for each route in parallel, speeding up the process
func HandleResponseDelayForRoute(responseChains []*tsdb.Chain, route parser.Routes, tsdbNameHash string, wg *sync.WaitGroup) {
	routeSuffix := filters.RouteDestroyer(route.URL)
	// Init paths for request-response-monitoring
	path := PathReqResDelay + "/" + "chunk_req_res_" + routeSuffix + ".json"
	c := make(chan Response)
	go RouteDispatcher(route, c)
	responseObject := <-c

	log.Printf("Writing responseObject to TSDB for %s", route.URL)

	block := *tsdb.GetNewBlock("req-res", getNormalizedBlockString(responseObject))

	for index := range responseChains {
		if responseChains[index].Path == path {
			responseChains[index] = responseChains[index].Append(block)
			responseChains[index].Save()
			break
		}
	}

	wg.Done()
}

// RouteDispatcher dispatches a route to respective handlers based on it's request type
func RouteDispatcher(route parser.Routes, c chan Response) {
	if route.Method == "GET" {
		res := HandleGetRequest(route.URL)
		c <- res
	} else {
		// Send a very large integer to automatically rule out as it
		// is much much larger than the threshold
		c <- Response{Delay: math.MaxInt32, ResLength: 0, ResStatusCode: 100}
	}
}

// HandleGetRequest specifically handles routes with GET Requests. Calculates timestamp before
// and after processing of each request and returns the difference
func HandleGetRequest(url string) Response {
	// Time init
	start := time.Now().UnixNano()
	resp := utils.SendGETRequest(url)
	resLength := resp.ContentLength
	respStatusCode := resp.StatusCode
	defer resp.Body.Close()

	end := time.Now().UnixNano()
	diff := int((end - start) / int64(time.Millisecond))

	return Response{Delay: diff, ResLength: resLength, ResStatusCode: respStatusCode}
}

// returns the stringified form of the combined data
func getNormalizedBlockString(b Response) string {
	return string(b.Delay) + tsdb.BlockDataSeparator + string(b.ResLength) + tsdb.BlockDataSeparator +
		string(b.ResStatusCode)
}
