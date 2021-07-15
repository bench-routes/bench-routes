package filters

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
