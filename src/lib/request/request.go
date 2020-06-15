package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Type-inputs for sending requests.
const (
	GET = iota
	POST
	PUT
	DELETE
	PATCH
)

// QuickInput implements the quick-input functionalities.
type QuickInput struct {
	headers, params, body map[string]string
	url                   string
}

// New returns a new quick-input type for implementing the
// quick-input testing functionality.
func New(url string, headers, params, body map[string]string) *QuickInput {
	return &QuickInput{
		url:     url,
		headers: headers,
		params:  params,
		body:    body,
	}
}

// Send sends the requests to the host/target. It can be executed
// parallely along with other goroutines.
func (q *QuickInput) Send(method uint, getResponse chan string) {
	switch method {
	case GET:
		if len(q.params) != 0 {
			q.url += "?" + q.formatParams()
		}
		var client http.Client
		request, err := http.NewRequest("GET", q.url, nil)
		if err != nil {
			panic(err)
		}
		q.formatHeaders(request)
		response, err := client.Do(request)
		defer response.Body.Close()
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		inStr := string(body)
		getResponse <- inStr
	}
}

// formatParams returns the params in the required format.
func (q *QuickInput) formatParams() string {
	if len(q.params) == 0 {
		return ""
	}
	var p string
	for k, v := range q.params {
		if k == "" {
			continue
		}
		p += fmt.Sprintf("%s=%s&", k, v)
	}
	if len(p) == 0 {
		return p
	}
	return p[0 : len(p)-1]
}

func (q *QuickInput) formatHeaders(request *http.Request) {
	for k, v := range q.headers {
		if k == "" {
			continue
		}
		request.Header.Set(k, v)
	}
}
