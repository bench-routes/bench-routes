package querier

import (
	"encoding/json"
	"fmt"
	"github.com/zairza-cetb/bench-routes/src/lib/logger"
	"time"
)

type response struct {
	Range        queryRange    `json:"range"`
	TimeInvolved time.Duration `json:"queryTime"`
	Value        interface{}   `json:"value"`
}

func getBaseResponse(r queryRange) (response, time.Time) {
	return response{
		Range: queryRange{
			Start: r.Start,
			End:   r.End,
		},
		Value: nil,
	}, time.Now()
}

func encode(r response) []byte {
	j, e := json.Marshal(r)
	if e != nil {
		logger.Terminal(fmt.Errorf("encoding error: %s", e.Error()).Error(), "p")
	}
	return j
}
