package evaluate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"unicode/utf8"
)

// Response contains evaluations of response delay and length, calculated during monitoring of an API.
type Response struct {
	Delay  time.Duration `json:"delay"`
	Length int           `json:"length"`
	Size   int           `json:"size"`
	Status int           `json:"status"`
}

// ExecuteMonitor monitors resDelay and resLength.
func Monitor(client *http.Client, request *http.Request) (*Response, error) {
	begin := time.Now()
	res, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error in sending request: %w", err)
	}
	resDelay := time.Since(begin)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error in reading response body: %w", err)
	}
	defer res.Body.Close()
	response := &Response{
		Delay:  resDelay,
		Length: utf8.RuneCountInString(string(resBody)),
		Size:   len(resBody),
		Status: res.StatusCode,
	}
	return response, nil
}
