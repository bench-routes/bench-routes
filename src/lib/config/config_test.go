package parser

import (
	"testing"
)

const (
	path = "../../../tests/configs/config-test.yml"
)

var (
	inst = Config{
		Address: path,
		Config:  &ConfigurationBR{},
	}
)

func TestLoad(t *testing.T) {
	inst = *inst.Load()
	res := *inst.Config
	if len(res.Interval) == 0 || len(res.Password) == 0 || len(res.Routes) == 0 {
		t.Errorf("faulty load of configuration.")
	}
}
