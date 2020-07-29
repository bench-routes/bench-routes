package parser

import (
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/prometheus/common/log"
)

// Config defines the structure type for implementing the interface
type Config struct {
	mutex   sync.RWMutex
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

// Body type for parameters passed along the url for specific route
type Body struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Route sets routes mentioned in configuration file
type Route struct {
	Method string    `yaml:"method"`
	URL    string    `yaml:"url"`
	Header []Headers `yaml:"headers"`
	Params []Params  `yaml:"params"`
	Body   []Body    `yaml:"body"`
	Labels []string  `yaml:"labels"`
}

// ResponseChangesConfig acts as a type for monitor-length configuration in config.yml
type ResponseChangesConfig struct {
	Mode float32 `yaml:"mode"`
	Mean float32 `yaml:"mean"`
}

// ServiceSignals type for defining current running states of various Services supported
// by BR. Allowed only two values: `active` OR `passive`
type ServiceSignals struct {
	Ping                  string `yaml:"ping"`
	FloodPing             string `yaml:"flood-ping"`
	Jitter                string `yaml:"jitter"`
	ReqResDelayMonitoring string `yaml:"req-res-delay-or-monitoring"`
}

// UConfig type for storing utilities in config.yml as local DB
type UConfig struct {
	RespChanges    ResponseChangesConfig `yaml:"monitor-length"`
	ServicesSignal ServiceSignals        `yaml:"services-state"`
}

// ConfigurationBR sets a type for configuration file which also acts as a local DB
type ConfigurationBR struct {
	Password  string     `yaml:"password"`
	UtilsConf UConfig    `yaml:"utils"`
	Interval  []Interval `yaml:"test_interval"`
	Routes    []Route    `yaml:"routes"`
}

// New returns an type for implementing the parser interface.
func New(path string) *Config {
	return &Config{
		Address: path,
	}
}

// Load loads the configuration file on startup.
func (inst *Config) Load() *Config {
	inst.mutex.RLock()
	defer inst.mutex.RUnlock()

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
	return inst
}

// Write force updates the configuration.
func (inst *Config) Write() (bool, error) {
	config := *inst.Config
	r, e := yaml.Marshal(config)
	if e != nil {
		log.Errorln(e.Error())
		return false, e
	}

	e = ioutil.WriteFile(inst.Address, r, 0644)
	if e != nil {
		panic(e)
	}

	return true, nil
}

// Refresh refreshes the Configuration settings.
func (inst *Config) Refresh() {
	inst.Load()
}

// AddRoute adds route to the Config.
func (inst *Config) AddRoute(route Route) {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	inst.Config.Routes = append(inst.Config.Routes, route)
	if _, err := inst.Write(); err != nil {
		panic(err)
	}
}

// TODO: Edit this to include labels
// GetNewRouteType returns a route based on the params provided.
func GetNewRouteType(method, url string, headers []Headers, params []Params, body []Body, labels []string) Route {
	return Route{
		Method: method,
		URL:    url,
		Header: headers,
		Params: params,
		Body:   body,
		Labels: labels,
	}
}
