package request

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/prometheus/common/log"
	config "github.com/zairza-cetb/bench-routes/src/lib/config"
)

// Type-inputs for sending requests.
const (
	GET = iota
	POST
	PUT
	DELETE
	PATCH
	// keep NULL to the bottom.
	NULL
)

// QuickInput implements the quick-input functionalities.
type QuickInput struct {
	headers, params, body map[string]string
	labels                []string
	url                   string
}

// ResponseWrapper wraps the response string with its status code.
type ResponseWrapper struct {
	Status             int
	ReponseStringified string
}

// New returns a new quick-input type for implementing the
// quick-input testing functionality.
func New(url string, headers, params, body map[string]string, labels []string) *QuickInput {
	return &QuickInput{
		url:     url,
		headers: headers,
		params:  params,
		body:    body,
		labels:  labels,
	}
}

// Send sends the requests to the host/target. It can be executed
// parallely along with other goroutines.
func (q *QuickInput) Send(method uint, getResponse chan ResponseWrapper) {
	var client http.Client
	switch method {
	case GET:
		if len(q.params) != 0 {
			q.url += "?" + q.formatParams()
		}
		request, err := http.NewRequest("GET", q.url, nil)
		if err != nil {
			panic(err)
		}
		q.applyHeaders(request)
		response, err := client.Do(request)
		if err != nil {
			log.Errorln(err)
		}
		if response == nil {
			fmt.Printf("nil response: %s\n", q.url)
			return
		}
		done(response.Body, response.StatusCode, getResponse)
		if err := response.Body.Close(); err != nil {
			fmt.Printf("error closing response body: %s\n", err.Error())
		}
	case POST:
		form := make(url.Values)
		q.populateBody(&form)
		request, err := http.NewRequest("POST", q.url, strings.NewReader(form.Encode()))
		if err != nil {
			panic(err)
		}
		request.PostForm = form
		q.applyHeaders(request)
		response, err := client.Do(request)
		if err != nil {
			log.Warnln(err)
		}
		if response == nil {
			log.Warnf("nil response: %s\n", q.url)
			return
		}
		done(response.Body, response.StatusCode, getResponse)
		if err := response.Body.Close(); err != nil {
			log.Warnf("error closing response body: %s\n", err.Error())
		}
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

func (q *QuickInput) applyHeaders(request *http.Request) {
	for k, v := range q.headers {
		if k == "" {
			continue
		}
		request.Header.Set(k, v)
	}
}

// populateBody applies the body as values to be assigned to the request.
func (q *QuickInput) populateBody(form *url.Values) {
	for k, v := range q.body {
		if k != "" || v != "" {
			form.Add(k, v)
		}
	}
}

func done(_body io.ReadCloser, status int, getResponse chan ResponseWrapper) {
	body, err := ioutil.ReadAll(_body)
	if err != nil {
		log.Warnln("err while reading body: ", err.Error())
		body = []byte("")
	}
	inStr := string(body)
	getResponse <- ResponseWrapper{ReponseStringified: inStr, Status: status}
}

// GetHeadersConfigFormatted converts the format of the headers
// into a config valid header format.
func (q *QuickInput) GetHeadersConfigFormatted() []config.Headers {
	var headers []config.Headers
	for k, v := range q.headers {
		headers = append(headers, config.Headers{OfType: k, Value: v})
	}
	return headers
}

// GetParamsConfigFormatted converts the format of the params
// into a config valid params format.
func (q *QuickInput) GetParamsConfigFormatted() []config.Params {
	var p []config.Params
	for k, v := range q.params {
		p = append(p, config.Params{Name: k, Value: v})
	}
	return p
}

// GetBodyConfigFormatted converts the format of the body
// into a config valid body format.
func (q *QuickInput) GetBodyConfigFormatted() []config.Body {
	var p []config.Body
	for k, v := range q.body {
		p = append(p, config.Body{Name: k, Value: v})
	}
	return p
}

// ToMap converts Headers, Params, Body slices to map.
func ToMap(slice interface{}) map[string]string {
	m := make(map[string]string)
	switch s := slice.(type) {
	case []config.Headers:
		for _, el := range s {
			m[el.OfType] = el.Value
		}
	case []config.Params:
		for _, el := range s {
			m[el.Name] = el.Value
		}
	case []config.Body:
		for _, el := range s {
			m[el.Name] = el.Value
		}
	}
	return m
}

// MethodUintPresentation takes http-method and returns the
// request compatible method input.
func MethodUintPresentation(method string) uint {
	switch method {
	case "GET":
		return GET
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "DELETE":
		return DELETE
	case "PATCH":
		return PATCH
	}
	return NULL
}
