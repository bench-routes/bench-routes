package utils

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// YAMLParser parses the yaml configuration files which are used as a local storage
type YAMLParser interface {
	Load() (bool, error)
	Write() (bool, error)
	Refresh() YAMLBenchRoutesType
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
