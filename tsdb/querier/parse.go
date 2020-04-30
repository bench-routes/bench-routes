package querier

import (
	"encoding/json"
	"fmt"
	"time"
)

type response struct {
	_range       queryRange    `json:"range"`
	timeInvolved time.Duration `json:"queryTime"`
	value        interface{}   `json:"value"`
}

type queryValue struct {
	timestamp int64 `yaml:"timestamp"`
	value     int64 `yaml:"value"`
}

const valueNULL string = "null"

func getBaseResponse(r queryRange) (response, time.Time) {
	return response{
		_range: queryRange{
			start: r.start,
			end:   r.end,
		},
		value: nil,
	}, time.Now()
}

func encode(r response) []byte {
	j, e := json.Marshal(r)
	if e != nil {
		fmt.Errorf("encoding error: %s", e.Error())
	}
	return j
}
