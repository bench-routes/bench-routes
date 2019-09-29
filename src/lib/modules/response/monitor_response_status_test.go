package response

import (
	"fmt"
	"testing"
)

func TestRequestForGETRoute(t *testing.T) {
	route := routes[0]
	res := HandleRequest(route)
	fmt.Println(res.status)
	if res.status == 200 {
		fmt.Println("Server found")
	} else {
		fmt.Println("server not found")
	}
}
