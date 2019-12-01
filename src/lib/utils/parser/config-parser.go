package parser

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// YAMLParser parses the yaml configuration files which are used as a local storage
type YAMLParser interface {
	Load() (bool, error)
	Write() (bool, error)
	Validate() bool
	Refresh() YAMLBenchRoutesType
}

// YAMLBenchRoutesType defines the structure type for implementing the interface
type YAMLBenchRoutesType struct {
	Address string
	Config  *ConfigurationBR
}

// Interval sets a type for intervals between consecutive similar tests
type Interval struct {
	Test     string `yaml:"test"`
	Type     string `yaml:"type"`
	Duration int64  `yaml:"duration"`
}

// Headers store the header values(ofType and value) from the config file
type Headers struct {
	OfType string `yaml:"type"`
	Value  string `yaml:"value"`
}

// Params type for parameters passed along the url for specific route
type Params struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Routes sets routes mentioned in configuration file
type Routes struct {
	Method string    `yaml:"method"`
	URL    string    `yaml:"url"`
	Route  string    `yaml:"route"`
	Header []Headers `yaml:"headers"`
	Params []Params  `yaml:"params"`
}

// ResponseChangesConfig acts as a type for response-length configuration in config.yml
type ResponseChangesConfig struct {
	Mode float32 `yaml:"mode"`
	Mean float32 `yaml:"mean"`
}

// ServiceSignals type for defining current running states of various services supported
// by BR. Allowed only two values: `active` OR `passive`
type ServiceSignals struct {
	Ping                  string `yaml:"ping"`
	FloodPing             string `yaml:"flood-ping"`
	Jitter                string `yaml:"jitter"`
	ReqResDelayMonitoring string `yaml:"req-res-delay-or-monitoring"`
}

// UConfig type for storing utilities in config.yml as local DB
type UConfig struct {
	RespChanges    ResponseChangesConfig `yaml:"response-length"`
	ServicesSignal ServiceSignals        `yaml:"services-state"`
}

// ConfigurationBR sets a type for configuration file which also acts as a local DB
type ConfigurationBR struct {
	Password  string     `yaml:"password"`
	Interval  []Interval `yaml:"test_interval"`
	Routes    []Routes   `yaml:"routes"`
	UtilsConf UConfig    `yaml:"utils"`
}

// Load loads the configuration file on startup
func (inst YAMLBenchRoutesType) Load() *YAMLBenchRoutesType {
	var yInstance ConfigurationBR
	file, e := ioutil.ReadFile(inst.Address)
	if e != nil {
		panic(e)
	}

	e = yaml.Unmarshal(file, &yInstance)
	if e != nil {
		panic(e)
	}
	inst.Config = &yInstance
	return &inst
}

func (inst YAMLBenchRoutesType) Write() (bool, error) {
	config := *inst.Config
	r, e := yaml.Marshal(config)
	if e != nil {
		log.Fatalf("%s\n", e)
		return false, e
	}

	e = ioutil.WriteFile(inst.Address, []byte(r), 0644)
	if e != nil {
		panic(e)
	}

	return true, nil
}

// Refresh refreshes the Configuration settings
func (inst YAMLBenchRoutesType) Refresh() YAMLBenchRoutesType {
	return *inst.Load()
}

// ValidateRoutesProp validates the `routes` property
// in the configuration file.
func ValidateRoutesProp(routes []Routes) {
	if len(routes) == 0 {
		log.Fatalf("`routes` property is missing.")
	} else {
		for i, route := range routes {
			if route.Method == "" {
				log.Fatalf("`method` property for route #%d is missing.\n", i)
			}
			if route.Route == "" {
				log.Fatalf("`route` property for route #%d is missing.\n", i)
			}
			if route.URL == "" {
				log.Fatalf("`url` property for route #%d is missing.\n", i)
			}
			if len(route.Header) != 0 {
				for i, header := range route.Header {
					if header.OfType == "" {
						log.Fatalf("`type` property for %s routes header #%d is missing.\n", route.Route, i+1)
					}
					if header.Value == "" {
						log.Fatalf("`value` property for %s routes header #%d is missing.\n", route.Route, i+1)
					}
				}
			}
			if len(route.Params) != 0 {
				for i, param := range route.Params {
					if param.Name == "" {
						log.Fatalf("`name` property for %s routes param #%d is missing.\n", route.Route, i+1)
					}
					if param.Value == "" {
						log.Fatalf("`value` property for %s routes param #%d is missing.\n", route.Route, i+1)
					}
				}
			}
		}
	}
}

// ValidateIntervalProp validates the `test_interval` property
// in the configuration file.
func ValidateIntervalProp(intervals []Interval) {
	if len(intervals) == 0 {
		log.Fatalf("`test_interval` property is missing.")
	} else {
		for i, interval := range intervals {
			if interval.Test == "" {
				log.Fatalf("`test` property for interval #%d is missing.\n", i)
			}
			if interval.Type == "" {
				log.Fatalf("`type` property for interval #%d is missing.\n", i)
			}
			if interval.Duration == 0 {
				log.Fatalf("`duration` property for interval #%d is missing.\n", i)
			}
		}
	}
}

// ValidateUtilsConf validates the `utils` property
// in the configuration file.
func ValidateUtilsConf(config *UConfig) {
	if config.ServicesSignal.FloodPing == "" {
		log.Fatalf("`flood-ping property` is missing.")
	}
	if config.ServicesSignal.Jitter == "" {
		log.Fatalf("`jitter` property is missing.")
	}
	if config.ServicesSignal.Ping == "" {
		log.Fatalf("`ping` property is missing.")
	}
	if config.ServicesSignal.ReqResDelayMonitoring == "" {
		log.Fatalf("`req-res-delay-or-monitoring` property is missing.")
	}
}

// Validate validates the config file.
func (inst YAMLBenchRoutesType) Validate() bool {
	var config = *inst.Config

	if len(config.Password) == 0 {
		log.Fatalf("`password` property is missing.")
	}

	ValidateRoutesProp(config.Routes)
	ValidateIntervalProp(config.Interval)
	ValidateUtilsConf(&config.UtilsConf)
	return true
}
