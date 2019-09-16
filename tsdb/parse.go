package tsdb

import (
	"io/ioutil"
	"errors"
)

func parse(path string) (*string, error) {
	res, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("file not existed. create a new chain")
	}
	str := string(res)
	return &str, nil
}