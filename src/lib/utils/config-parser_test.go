package utils

import (
	"testing"
)

const (
	path = "../../../test-files/configs/config-test.yml"
)

var (
	inst = YAMLBenchRoutesType{
		address: path,
		config:  &ConfigurationBR{},
	}
)

func TestLoad(t *testing.T) {
	inst = *inst.Load()
	res := *inst.config
	if len(res.Interval) == 0 || len(res.Password) == 0 || len(res.Routes) == 0 {
		t.Errorf("faulty load of configuration.")
	} else {
		t.Log(inst)
	}
}
