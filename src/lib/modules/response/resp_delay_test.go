package response

import (
	"log"
	"math"
	"testing"
)

const (
	threshold = 0
)

var routes = []Route{
	Route{
		"https://www.zairza.in/",
		"GET",
	},
	Route{
		"http://www.zairza.in",
		"OPTIONS",
	},
}

// what happens when requests other than GET, POST, PUT, DELETE hop in
func TestRouteDispatcherForUnmentionedRequestTypes(t *testing.T) {
	routeToTest := routes[1]
	res := HandleResponseDelayForRoute(routeToTest)
	// Output: Should return a large Integer
	if res.delay != math.MaxInt32 {
		t.Errorf("should return a large integer to automatically rule out of threshold")
	}
	log.Println(res)
}

// Test if routes that have larger delays in response times
// are conveyed to the user via the alerter
func TestIfLargerReponseDelaysAreProneToAlert(t *testing.T) {
	routeToTest := routes[0]
	res := HandleResponseDelayForRoute(routeToTest)
	// Output: res > threshold, Prone to alerting
	if res.delay < threshold {
		t.Errorf("Invalid result. Expected result > threshold")
	}
	log.Println(res)
}

// Test if Content-Length headers are not set, resLength
// returns 0 or any negative number
func TestIfNoContentLengthHeadersReturnsNegative(t *testing.T) {
	routeToTest := routes[1]
	res := HandleResponseDelayForRoute(routeToTest)
	// Output: 0 or a negative number
	if res.resLength > 0 {
		t.Errorf("Should return a 0 or a negative number as Content-Length header is not set!")
	}
	log.Println(res)
}
