package configparser

import (
	"testing"
)

const (
	path = "../../../tests/configs/config-test_v2.yml"
)

var (
	inst = Config{
		Address: path,
		Root:    &RootConfig{},
	}
)

func TestLoad(t *testing.T) {
	config, err := inst.Load()
	if err != nil {
		t.Error(err)
	}
	err = config.Validate()
	if err != nil {
		t.Error(err)
	}
}
