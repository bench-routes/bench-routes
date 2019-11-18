package tsdb

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"
	"strings"
)

func parse(path string) (*string, error) {
	res, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("file not existed. create a new chain")
	}
	str := string(res)
	return &str, nil
}

// Parser type for accessing parsing functions
type Parser struct{}

// Parse acts as an global interface for converting different types into the required structure
type Parse interface {
	ParseToJSON(a []Block) []byte
}

// ParseToJSON converts the chain into Marshallable JSON
func (p Parser) ParseToJSON(a []Block) (j []byte) {
	b := []BlockJSON{}

	for _, inst := range a {
		t := BlockJSON{
			Timestamp:      inst.Timestamp,
			NormalizedTime: inst.NormalizedTime,
			Datapoint:      inst.Datapoint,
		}
		b = append(b, t)
	}

	j, e := json.Marshal(b)
	if e != nil {
		panic(e)
	}
	return
}

func loadFromStorage(raw *string) *[]BlockJSON {
	inst := []BlockJSON{}
	b := []byte(*raw)
	e := json.Unmarshal(b, &inst)
	if e != nil {
		panic(e)
	}
	return &inst
}

func formLinkedChainFromRawBlock(a *[]BlockJSON) *[]Block {
	r := *a
	l := len(r)
	arr := []Block{}
	for i := 0; i < l; i++ {
		inst := Block{
			Timestamp:      r[i].Timestamp,
			NormalizedTime: r[i].NormalizedTime,
			Datapoint:      r[i].Datapoint,
		}
		arr = append(arr, inst)
	}

	return &arr
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
func GetTimeStamp() string {
	t := time.Now()

	return s(t.Year()) + "|" + s(t.Month()) + "|" + s(t.Day()) + "|" + s(t.Hour()) + "|" +
			s(t.Minute()) + "|" + s(t.Second()) + "#" + s(milliSeconds())
}

// GetNormalizedTime returns the UnixNano time as normalized time.
func GetNormalizedTime() int64 {
	return time.Now().UnixNano()
}

func s(st interface{}) string {
	return st.(string)
}

func milliSeconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
