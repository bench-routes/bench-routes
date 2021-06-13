package config

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

//Config struct stores the APIs parsed from local-config.yml file.
type Config struct {
	path string
	mux  sync.RWMutex
	APIs []API `yaml:"apis"`
}

//API stores the information of the endpoints.
type API struct {
	Name    string            `yaml:"name,omitempty"`
	Every   time.Duration     `yaml:"every,omitempty"`
	Domain  string            `yaml:"domain_or_ip,omitempty"`
	Route   string            `yaml:"route,omitempty"`
	Method  string            `yaml:"method,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
	Params  map[string]string `yaml:"params,omitempty"`
	Body    map[string]string `yaml:"body,omitempty"`
}

//Instantiate a Config object and reloads the data from the given file location.
func NewConf(path string) (*Config, error) {
	c := &Config{
		path: path,
	}
	c, err := c.Reload()
	if err != nil {
		return nil, fmt.Errorf("creating config: %w", err)
	}
	return c, nil
}

//Reloads data from the config file.
func (c *Config) Reload() (*Config, error) {
	conf := new(Config)

	c.mux.RLock()
	file, err := ioutil.ReadFile(c.path)
	c.mux.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	if err = yaml.Unmarshal(file, conf); err != nil {
		return nil, fmt.Errorf("unmarshalling file content: %w", err)
	}

	return conf, nil
}

//Adds API to the Config struct
func (c *Config) AddAPI(api API) (*Config, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.APIs = append(c.APIs, api)
	return c, nil
}
