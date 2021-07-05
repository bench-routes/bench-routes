package evaluate

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

// Ping icmp requests to target.
type Ping struct {
	Min  time.Duration `json:"min"`
	Mean time.Duration `json:"mean"`
	Max  time.Duration `json:"max"`
}

// Jitter is the mean variation between consecutive ping operations.
type Jitter struct {
	Value time.Duration `json:"jitter"`
}

// Machine evaluates calculations related to ping and jitter.
func Machine(host string) (*Ping, *Jitter, error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating pinger : %w", err)
	}
	pinger.Count = 5
	var lastTime, sum time.Duration
	// Calculating jitter using ping values.
	pinger.OnRecv = func(pkt *ping.Packet) {
		if lastTime != time.Second*0 {
			sum += absDiff(lastTime, pkt.Rtt)
		}
		lastTime = pkt.Rtt
	}
	// Runing the pinger.
	if err := pinger.Run(); err != nil {
		return nil, nil, fmt.Errorf("error running ping : %w", err)
	}
	stats := pinger.Statistics()
	jitter := sum / time.Duration(pinger.Count-1)
	pingVal := &Ping{
		Min:  stats.MaxRtt,
		Mean: stats.AvgRtt,
		Max:  stats.MinRtt,
	}
	jitterVal := &Jitter{
		Value: jitter,
	}

	return pingVal, jitterVal, nil
}

func absDiff(a, b time.Duration) time.Duration {
	if a >= b {
		return a - b
	}
	return b - a
}
