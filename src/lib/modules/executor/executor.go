package executor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bench-routes/bench-routes/tsdb/file"
	"github.com/go-ping/ping"
)

type Ping struct {
	Min time.Duration `json:"minVal"`
	Max time.Duration `json:"maxVal"`
	Avg time.Duration `json:"avgVal"`
}

type Jitter struct {
	Min time.Duration `json:"minVal"`
	Max time.Duration `json:"maxVal"`
	Avg time.Duration `json:"avgVal"`
}

type Monitor struct {
	ResDelay  time.Duration `json:"resDelay"`
	ResLength int64         `json:"resLength"`
}

// ExecuteMachine executes ping and jitter operations
func ExecuteMachine(app file.Appendable, host string) error {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return fmt.Errorf("error creating ping : %w", err)
	}
	pinger.Count = 5
	var lastTime time.Duration
	var sum time.Duration
	var minJitter time.Duration = time.Hour * 24
	var maxJitter time.Duration = 0
	// Calculating jitter using ping values
	pinger.OnRecv = func(pkt *ping.Packet) {
		if lastTime != time.Second*0 {
			minJitter = min(minJitter, absDiff(lastTime, pkt.Rtt))
			maxJitter = max(minJitter, absDiff(lastTime, pkt.Rtt))
			sum += absDiff(lastTime, pkt.Rtt)
		}
		lastTime = pkt.Rtt
	}
	// Runing the pinger
	if err := pinger.Run(); err != nil {
		return fmt.Errorf("error running ping : %w", err)
	}
	stats := pinger.Statistics()
	avgjitter := sum / time.Duration(pinger.Count-1)
	pingval := fmt.Sprintf("%v|%v|%v", stats.MaxRtt.Nanoseconds(), stats.MinRtt.Nanoseconds(), stats.AvgRtt.Nanoseconds())
	jitterval := fmt.Sprintf("%v|%v|%v", maxJitter.Nanoseconds(), minJitter.Nanoseconds(), avgjitter.Nanoseconds())

	fmt.Println("Ping :", pingval)
	fmt.Println("Jitter :", jitterval)
	app.Append(file.NewBlock("ping", pingval))
	app.Append(file.NewBlock("jitter", jitterval))
	return nil
}

func absDiff(a, b time.Duration) time.Duration {
	if a >= b {
		return a - b
	}
	return b - a
}

func max(a, b time.Duration) time.Duration {
	if a >= b {
		return a
	}
	return b
}

func min(a, b time.Duration) time.Duration {
	if a >= b {
		return b
	}
	return a
}

// ExecuteMonitor monitors resDelay and resLength
func ExecuteMonitor(app file.Appendable, client *http.Client, request *http.Request) error {
	stamp := time.Now()
	res, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error in sending request: %w", err)
	}
	resDelay := time.Since(stamp)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error in reading response body: %w", err)
	}
	res.Body.Close()
	val := fmt.Sprintf("%v|%v", resDelay.Nanoseconds(), len(resBody))
	fmt.Println(val)
	app.Append(file.NewBlock("monitoring", val))
	return nil
}
