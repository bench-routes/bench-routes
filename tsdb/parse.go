package tsdb

import (
	"io/ioutil"
	"errors"
	"encoding/json"
	"fmt"
)

func parse(path string) (*string, error) {
	res, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("file not existed. create a new chain")
	}
	str := string(res)
	return &str, nil
}

func loadFromStorage(raw *string) *[]BlockJSON {
	inst := []BlockJSON{}
	b := []byte(*raw)
	e := json.Unmarshal(b, &inst)
	fmt.Println(inst)
	if e != nil {
		panic(e)
	}
	return &inst
}