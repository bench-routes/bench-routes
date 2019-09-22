package response

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

// Route struct
type Route struct {
	url         string
	requestType string
}

// HandleResponseDelayForRoute is the initial entrypoint function for this module which takes
// in a Route struct and supplies it to a function in turn to handle it accordingly. We create
// channels to run tests for each route in parallel, speeding up the process
func HandleResponseDelayForRoute(route Route) int {
	c := make(chan int)
	go RouteDispatcher(route, c)
	timeElapsed := <-c
	return timeElapsed
}

// RouteDispatcher dispatches a route to respective handlers based on it's request type
func RouteDispatcher(route Route, c chan int) {
	if route.requestType == "GET" {
		respTime := HandleGetRequest(route.url)
		c <- respTime
	} else {
		// Send a very large integer to automatically rule out as it
		// is much much larger than the threshold
		c <- math.MaxInt32
	}
}

// HandleGetRequest specifically handles routes with GET Requests. Calculates timestamp before
// and after processing of each request and returns the difference
func HandleGetRequest(url string) int {
	// Time init
	start := time.Now().UnixNano()

	resp, err := http.Get(url)
	if err != nil {
		// Prone to alerting, printing for now
		fmt.Println(err)
	}
	defer resp.Body.Close()

	end := time.Now().UnixNano()
	diff := int((end - start) / int64(time.Millisecond))

	return diff
}
