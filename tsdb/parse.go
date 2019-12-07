package tsdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func parse(path string) (*string, error) {
	res, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("file not existed. create a new chain")
	}
	str := string(res)
	return &str, nil
}

// parseToJSON converts the chain into Marshallable JSON
func parseToJSON(a []Block) (j []byte) {
	j, e := json.Marshal(a)
	if e != nil {
		panic(e)
	}
	return
}

func loadFromStorage(raw *string) *[]Block {
	inst := []Block{}
	b := []byte(*raw)
	e := json.Unmarshal(b, &inst)
	if e != nil {
		panic(e)
	}
	return &inst
}

func checkAndCreatePath(path string) {
	array := strings.Split(path, "/")
	array = array[:len(array)-1]
	path = strings.Join(array, "/")
	_, err := os.Stat(path)
	if err != nil {
		e := os.MkdirAll(path, os.ModePerm)
		if e != nil {
			panic(e)
		}
	}
}

func saveToHDD(path string, data []byte) error {
	checkAndCreatePath(path)
	e := ioutil.WriteFile(path, data, 0755)
	if e != nil {
		return e
	}
	return nil
}

// GetTimeStamp returns the timestamp
func GetTimeStampCalc() string {
	t := time.Now()

	return s(t.Year()) + "|" + s(t.Month()) + "|" + s(t.Day()) + "|" + s(t.Hour()) + "|" +
		s(t.Minute()) + "|" + s(t.Second()) + "#" + s(milliSeconds())
}

// GetNormalizedTime returns the UnixNano time as normalized time.
func GetNormalizedTimeCalc() int64 {
	return time.Now().UnixNano()
}

func s(st interface{}) string {
	return fmt.Sprintf("%v", st)
}

func milliSeconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
