package response

import (
	"math"
	"time"

	"github.com/zairza-cetb/bench-routes/src/lib/utils"
)

// Route struct
type Route struct {
	url         string
	requestType string
}

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
func HandleResponseDelayForRoute(route Route) Response {
	c := make(chan Response)
	go RouteDispatcher(route, c)
	responseObject := <-c
	return responseObject
}

// RouteDispatcher dispatches a route to respective handlers based on it's request type
func RouteDispatcher(route Route, c chan Response) {
	if route.requestType == "GET" {
		res := HandleGetRequest(route.url)
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
