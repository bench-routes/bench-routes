package querier

import (
	"encoding/json"
	"fmt"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"time"
)

// QueryResponse is the response sent after processing the query.
type QueryResponse struct {
	Range        queryRange    `json:"range"`
	TimeInvolved time.Duration `json:"queryTime"`
	Value        interface{}   `json:"values"`
}

type queryValue struct {
	Value          interface{} `json:"value"`
	Timestamp      interface{} `json:"timestamp"`
	NormalizedTime int64       `json:"normalizedTime"`
}

func getBaseResponse(r queryRange) (QueryResponse, time.Time) {
	return QueryResponse{
		Range: queryRange{
			Start: r.Start,
			End:   r.End,
		},
		Value: nil,
	}, time.Now()
}

func encode(r QueryResponse) []byte {
	j, e := json.Marshal(r)
	if e != nil {
		logger.Terminal(fmt.Errorf("encoding error: %s", e.Error()).Error(), "p")
	}
	return j
}
