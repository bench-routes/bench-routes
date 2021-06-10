package configparser

import (
	"io/ioutil"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Address string
	Mutex   sync.RWMutex
	Root    *RootConfig
}

type RootConfig struct {
	APIs []API `yaml:"apis"`
}

type API struct {
	Name    string            `yaml:"name"`
	Every   time.Duration     `yaml:"every"`
	Domain  string            `yaml:"domain_or_ip"`
	Route   string            `yaml:"route"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Params  map[string]string `yaml:"params"`
	Body    map[string]string `yaml:"body"`
}

func NewConf(configPath string) *Config {
	return &Config{
		Address: configPath,
	}
}

func (c *Config) Reload() (*Config, error) {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	var confInstance RootConfig

	file, err := ioutil.ReadFile(c.Address)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &confInstance)

	if err != nil {
		return nil, err
	}
	c.Root = &confInstance
	return c, nil
}

func (c *Config) WriteConf() error {
	root := *c.Root
	file, err := yaml.Marshal(root)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.Address, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) AddAPI(api API) (*Config, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Root.APIs = append(c.Root.APIs, api)
	err := c.WriteConf()
	if err != nil {
		return nil, err
	}
	return c, nil
}
