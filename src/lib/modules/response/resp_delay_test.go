package response

// var routes = []utils.Routes{
// 	utils.Routes{
// 		Method: "GET",
// 		URL:    "https://www.zairza.in/",
// 		Route:  "/test",
// 	},
// 	utils.Routes{
// 		Method: "OPTIONS",
// 		URL:    "https://www.zairza.in/",
// 		Route:  "/test",
// 	},
// }

// func initTest()

// // what happens when requests other than GET, POST, PUT, DELETE hop in
// func TestRouteDispatcherForUnmentionedRequestTypes(t *testing.T) {
// 	routeToTest := routes[1]
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go HandleResponseDelayForRoute(routeToTest, utils.GetHash(routeToTest.URL), &wg)
// 	wg.Wait()
// }

// // Test if routes that have larger delays in response times
// // are conveyed to the user via the alerter
// func TestIfLargerReponseDelaysAreProneToAlert(t *testing.T) {
// 	routeToTest := routes[0]
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go HandleResponseDelayForRoute(routeToTest, utils.GetHash(routeToTest.URL), &wg)
// 	wg.Wait()
// }

// // Test if Content-Length headers are not set, resLength
// // returns 0 or any negative number
// func TestIfNoContentLengthHeadersReturnsNegative(t *testing.T) {
// 	routeToTest := routes[1]
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go HandleResponseDelayForRoute(routeToTest, utils.GetHash(routeToTest.URL), &wg)
// 	wg.Wait()
// }
