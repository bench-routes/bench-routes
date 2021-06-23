package evaluate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-ping/ping"
)

// Ping icmp requests to target.
type Ping struct {
	Min	 	time.Duration 	`json:"min"`
	Mean 	time.Duration 	`json:"mean"`
	Max 	time.Duration 	`json:"max"`
}

// Jitter is the mean variation between consecutive ping operations.
type Jitter struct {
	Value 	time.Duration 	`json:"jitter"`
}

// Response contains evaluations of response delay and length, calculated during monitoring of an API.
type Response struct {
	Delay  	time.Duration 	`json:"delay"`
	Length 	int		    	`json:"length"`
	Size  	[]byte       	`json:"size"`
}

// Machine evaluates calculations related to ping and jitter.
func Machine(host string) (*Ping,*Jitter,error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return nil,nil,fmt.Errorf("error creating ping : %w", err)
	}
	pinger.Count = 5
	var lastTime time.Duration
	var sum time.Duration
	// Calculating jitter using ping values.
	pinger.OnRecv = func(pkt *ping.Packet) {
		if lastTime != time.Second*0 {
			sum += absDiff(lastTime, pkt.Rtt)
		}
		lastTime = pkt.Rtt
	}
	// Runing the pinger
	if err := pinger.Run(); err != nil {
		return nil,nil,fmt.Errorf("error running ping : %w", err)
	}
	stats := pinger.Statistics()
	avgjitter := sum / time.Duration(pinger.Count-1)
	pingval := &Ping{
		Min: 	stats.MaxRtt,
		Mean:	stats.AvgRtt,
		Max: 	stats.MinRtt,
	}
	jitterval := &Jitter{
		Value: avgjitter,
	}
	
	return pingval,jitterval,nil
}

func absDiff(a, b time.Duration) time.Duration {
	if a >= b {
		return a - b
	}
	return b - a
}

// ExecuteMonitor monitors resDelay and resLength
func Monitor(client *http.Client, request *http.Request) (*Response,error) {
	stamp := time.Now()
	res, err := client.Do(request)
	if err != nil {
		return nil,fmt.Errorf("error in sending request: %w", err)
	}
	resDelay := time.Since(stamp)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil,fmt.Errorf("error in reading response body: %w", err)
	}
	res.Body.Close()
	response := &Response{
		Delay: 	resDelay,
		Length: len(resBody),
		Size: 	resBody,
	}
	return response,nil
}
