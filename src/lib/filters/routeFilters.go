package filters

import (
	"encoding/json"

	"github.com/zairza-cetb/bench-routes/src/lib/config"
)

// Params JSON implementation of config-parser.go
type Params struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Headers JSON implementation of config-parser.go
type Headers struct {
	OfType string `json:"type"`
	Value  string `json:"value"`
}

// Routes defines a struct for websocket communication inorder to be Marshalable
type Routes struct {
	Method string    `json:"method"`
	URL    string    `json:"url"`
	Route  string    `json:"route"`
	Header []Headers `json:"headers"`
	Params []Params  `json:"params"`
}

// RouteJSONParser acts as a JSON parser for the Routes in the config-parser.go
type RouteJSONParser struct {
	Routes []Routes `json:"routes"`
}

func RouteYAMLtoJSONParser(r []parser.Route) (rr []byte) {
	var rp RouteJSONParser
	for _, route := range r {
		var tmpH []Headers
		var tmpP []Params
		for _, head := range route.Header {
			tmp := Headers{
				OfType: head.OfType,
				Value:  head.Value,
			}
			tmpH = append(tmpH, tmp)
		}

		for _, param := range route.Params {
			tmp := Params{
				Name:  param.Name,
				Value: param.Value,
			}
			tmpP = append(tmpP, tmp)
		}
		tmp := Routes{
			Method: route.Method,
			URL:    route.URL,
			Header: tmpH,
			Params: tmpP,
		}
		rp.Routes = append(rp.Routes, tmp)
	}
	rr, err := json.Marshal(rp)
	if err != nil {
		panic(err)
	}
	return rr
}
