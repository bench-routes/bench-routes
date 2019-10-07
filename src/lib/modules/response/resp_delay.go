package response

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/utils"
)

// Path to storage in tsdb for Req res delay
const (
	// PathPing stores the defualt address of storage directory of ping data
	PathReqResDelay = "storage/req-res-delay"
)

// Response struct
// This is the object that we return from this module
// Contains delay in response and the response length
type Response struct {
	delay         int
	resLength     int64
	resStatusCode int
}

// HandleResponseDelayForRoute is the initial entrypoint function for this module which takes
// in a Route struct and supplies it to a function in turn to handle it accordingly. We create
// channels to run tests for each route in parallel, speeding up the process
func HandleResponseDelayForRoute(route utils.Routes, tsdbNameHash string, wg *sync.WaitGroup) {
	tsdbNameHash = PathReqResDelay + "/" + "chunk_req_res_" + tsdbNameHash + ".json"
	c := make(chan Response)
	go RouteDispatcher(route, c)
	responseObject := <-c
	// TODO: Do not print the responseObject but
	// store it in the tsdb as a record
	fmt.Println(responseObject)
	wg.Done()
}

// RouteDispatcher dispatches a route to respective handlers based on it's request type
func RouteDispatcher(route utils.Routes, c chan Response) {
	if route.Method == "GET" {
		res := HandleGetRequest(route.URL)
		c <- res
	} else {
		// Send a very large integer to automatically rule out as it
		// is much much larger than the threshold
		c <- Response{delay: math.MaxInt32, resLength: 0, resStatusCode: 100}
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

	return Response{delay: diff, resLength: resLength, resStatusCode: respStatusCode}
}
