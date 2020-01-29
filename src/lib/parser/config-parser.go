package parser

import (
	"io/ioutil"

	"github.com/zairza-cetb/bench-routes/src/lib/logger"
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
	Duration *int64 `yaml:"duration"`
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
	Mode *float32 `yaml:"mode"`
	Mean *float32 `yaml:"mean"`
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

// New returns an type for implementing the parser interface.
func New(path string) YAMLBenchRoutesType {
	return YAMLBenchRoutesType{
		Address: path,
	}
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
		logger.Terminal(e.Error(), "f")
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
