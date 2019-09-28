package response

import(
	"net/http"
	"github.com/zairza-cetb/bench-routes/src/lib/utils"
)

//ResponseStatus contains the status code of the requested url
type ResponseStatus struct {
	status		int
}


//GETRequestDispatcher dispatches a specific type(GET,POST,PUT,DETLETE,)route to the respective function
func GETRequestDispatcher(url string, chnl chan *http.Response){
	res:=utils.SendGETRequest(url)
	chnl<-res
}
//HandleRequest is the entry point for this module
func HandleRequest(route Route) ResponseStatus{
	chnl := make(chan *http.Response)
	if(route.requestType=="GET"){
		go GETRequestDispatcher(route.url,chnl)
		resp:=<-chnl
		return ResponseStatus{status:resp.StatusCode}
	} 
		return ResponseStatus{status:100}
	
}

