package configparser

import (
	"fmt"
	"log"
	"testing"
)

const (
	path = "../../../tests/configs/config-test_v2.yml"
)

var (
	inst = Config{
		Address: path,
		Root: &RootConfig{},
	}
)

func TestLoad(t *testing.T) {
	config,err := inst.Load()
	if err != nil {
		log.Fatal(err)
		return 
	}
	res := *config.Root
	// if len(res.Interval) == 0 || len(res.Password) == 0 || len(res.Routes) == 0 {
	// 	t.Errorf("faulty load of configuration.")
	// }
	fmt.Println(res)
}