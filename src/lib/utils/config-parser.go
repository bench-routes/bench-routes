package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// YAMLParser parses the yaml configuration files which are used as a local storage
type YAMLParser interface {
	Load() (bool, error)
	Write() (bool, error)
}

// YAMLBenchRoutesType defines the structure type for implementing the interface
type YAMLBenchRoutesType struct {
	address string
	config  *ConfigurationBR
}

// Interval sets a type for intervals between consecutive similar tests
type Interval struct {
	Test     string `yaml:"test"`
	Type     string `yaml:"type"`
	Duration int64  `yaml:"duration"`
}

type headers struct {
	OfType string `yaml:"type"`
	Value  string `yaml:"value"`
}

type params struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Routes sets routes mentioned in configuration file
type Routes struct {
	Method string    `yaml:"method"`
	URL    string    `yaml:"url"`
	Route  string    `yaml:"route"`
	Header []headers `yaml:"headers"`
	Params []params  `yaml:"params"`
}

// ResponseChangesConfig acts as a type for response-length configuration in config.yml
type ResponseChangesConfig struct {
	Mode float32 `yaml:"mode"`
	Mean float32 `yaml:"mean"`
}

// UConfig type for storing utilities in config.yml as local DB
type UConfig struct {
	RespChanges ResponseChangesConfig `yaml:"response-length"`
}

// ConfigurationBR sets a type for configuration file which also acts as a local DB
type ConfigurationBR struct {
	Password  string     `yaml:"password"`
	Interval  []Interval `yaml:"interval"`
	Routes    []Routes   `yaml:"routes"`
	UtilsConf UConfig    `yaml:"utils"`
}

// Load loads the configuration file on startup
func (inst YAMLBenchRoutesType) Load() *YAMLBenchRoutesType {
	var yInstance ConfigurationBR
	file, e := ioutil.ReadFile(inst.address)
	if e != nil {
		panic(e)
	}

	e = yaml.Unmarshal(file, &yInstance)
	if e != nil {
		panic(e)
	}
	inst.config = &yInstance
	log.Println(yInstance)
	return &inst
}

func (inst YAMLBenchRoutesType) Write() (bool, error) {
	config := *inst.config
	r, e := yaml.Marshal(config)
	if e != nil {
		log.Fatalf("%s\n", e)
		return false, e
	}

	e = ioutil.WriteFile(inst.address, []byte(r), 0644)
	if e != nil {
		panic(e)
	}

	return true, nil
}
