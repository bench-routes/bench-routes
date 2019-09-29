package response

import (
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
	"net/http"
)

//Status contains the status code of the requested url
type Status struct {
	status int
}

//GETRequestDispatcher dispatches a specific type(GET,POST,PUT,DETLETE,)route to the respective function
func GETRequestDispatcher(url string, chnl chan *http.Response) {
	res := utils.SendGETRequest(url)
	chnl <- res
}

//HandleRequest is the entry point for this module
func HandleRequest(route Route) Status {
	chnl := make(chan *http.Response)
	if route.requestType == "GET" {
		go GETRequestDispatcher(route.url, chnl)
		resp := <-chnl
		return Status{status: resp.StatusCode}
	}
	return Status{status: 100}

}
