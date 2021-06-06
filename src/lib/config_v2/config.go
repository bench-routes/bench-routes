package configparser

import (
	"io/ioutil"
	"sync"
	"time"

	"github.com/prometheus/common/log"
	//"gopkg.in/yaml.v2"
	"gopkg.in/yaml.v3"
)

type Config struct{
	Address 	string
	Mutex 		sync.RWMutex 
	Root 		*RootConfig
}

type RootConfig struct{
	APIs []API 	`yaml:"apis"`
}

type API struct{
	Name 		string 				`yaml:"name"`
	Every 		time.Duration 		`yaml:"every"`
	Domain 		string				`yaml:"domian_or_ip"`
	Route 		string				`yaml:"route"`
	Method 		string				`yaml:"method"`
	Headers 	map[string]string	`yaml:"headers"`
	Params		map[string]string	`yaml:"params"`
	Body		map[string]string	`yaml:"body"`	
}


func NewConf(configPath string) (*Config){
	return &Config{
		Address: configPath,
	}	
}

func(c *Config) Load()(*Config,error){
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	var confInstance RootConfig

	file,err := ioutil.ReadFile(c.Address)
	if err != nil {
		log.Error("Error reading file")
		return nil,err
	}

	err = yaml.Unmarshal(file,&confInstance)

	if err != nil {
		log.Error("Error marshalling config file")
		return nil,err;
	}
	c.Root = &confInstance
	return c,nil
}
