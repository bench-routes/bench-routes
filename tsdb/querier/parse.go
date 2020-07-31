package querier

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"math"
	"time"
)

// QueryResponse is the response sent after processing the query.
type QueryResponse struct {
	TimeSeriesPath string        `json:"timeSeriesPath"`
	Range          queryRange    `json:"range"`
	TimeInvolved   time.Duration `json:"queryTime"`
	Value          interface{}   `json:"values"`
}

type queryValue struct {
	Value          interface{} `json:"value"`
	Timestamp      interface{} `json:"timestamp"`
	NormalizedTime int64       `json:"normalizedTime"`
}

func getBaseResponse(r *queryRange) (QueryResponse, time.Time) {
	if r == nil {
		r = &queryRange{
			Start: math.MaxInt64,
			End:   math.MinInt64,
		}
	}
	return QueryResponse{
		Range: queryRange{
			Start: r.Start,
			End:   r.End,
		},
		Value: nil,
	}, time.Now()
}

func encode(r QueryResponse, enabled bool) interface{} {
	if enabled {
		j, e := json.Marshal(r)
		if e != nil {
			log.Errorln(fmt.Errorf("encoding error: %s", e.Error()).Error())
		}
		return j
	}
	return r
}
