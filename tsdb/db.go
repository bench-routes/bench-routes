package tsdb

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/prometheus/common/log"
)

const (
	// BlockDataSeparator sets a separator for block data value.
	BlockDataSeparator = "|"
	// ParameterSeparator sets a separator between parameters in a block.
	ParameterSeparator = ","
	// FileExtension file extension for db files.
	FileExtension = ".json"
)

//HashList use case HashList for the HashTable
type HashList struct {
	URL        string `json:"url"`         // stores the URL
	HashNumber string `json:"hash-number"` // stores the hash number of a particular URL and type
	Type       string `json:"type"`        // type of action
}

func init() {
	if _, err := parse("storage/000000"); err != nil {
		log.Infof("creating in-memory chain2: storage/0000000\n")
		if err := saveToHDD("storage/000000", []byte("[]")); err != nil {
			panic(err)
		}
	}
}

func parse(path string) (*string, error) {
	res, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("file not existed. create a new chain")
	}
	str := string(res)
	return &str, nil
}

func checkAndCreatePath(path string) {
	array := strings.Split(path, "/")
	array = array[:len(array)-1]
	path = strings.Join(array, "/")
	if _, err := os.Stat(path); err != nil {
		if e := os.MkdirAll(path, os.ModePerm); e != nil {
			panic(e)
		}
	}
}

func saveToHDD(path string, data []byte) error {
	checkAndCreatePath(path)
	if e := ioutil.WriteFile(path, data, 0755); e != nil {
		return e
	}
	return nil
}

//SaveIndexTable saves the index table to HDD
func SaveIndexTable(path string, data *map[string]*HashList) error {
	checkAndCreatePath(path)
	res := formatToStoreData(data)
	write := []byte(res)
	if e := ioutil.WriteFile(path, write, 0755); e != nil {
		return e
	}
	return nil
}

func formatToStoreData(data *map[string]*HashList) string {
	var result string = ""
	for _, element := range *data {
		result += element.HashNumber + ParameterSeparator + element.Type + ParameterSeparator + element.URL + BlockDataSeparator
	}
	return result
}
