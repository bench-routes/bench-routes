package tsdb

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
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

	// form doubly linked list
	for i := 0; i < l; i++ {
		if i == 0 {
			arr[i].PrevBlock = nil
		} else {
			arr[i].PrevBlock = &arr[i-1]
		}
		if i == l-1 {
			arr[i].NextBlock = nil
		} else {
			arr[i].NextBlock = &arr[i+1]
		}
	}
	return &arr
}

func loadFromStoragePing(raw *string) *[]BlockPingJSON {
	inst := []BlockPingJSON{}
	b := []byte(*raw)
	e := json.Unmarshal(b, &inst)
	if e != nil {
		panic(e)
	}
	return &inst
}

func loadFromStorageFloodPing(raw *string) *[]BlockFloodPingJSON {
	inst := []BlockFloodPingJSON{}
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

// ParseToJSONPing converts the ping chain into Marshallable JSON
func (p Parser) ParseToJSONPing(a []BlockPing) (j []byte) {
	b := []BlockPingJSON{}
	for _, inst := range a {
		t := BlockPingJSON{
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

// ParseToJSONFloodPing converts the flood ping chain into Marshallable JSON
func (p Parser) ParseToJSONFloodPing(a []BlockFloodPing) (j []byte) {
	b := []BlockFloodPingJSON{}
	for _, inst := range a {
		t := BlockFloodPingJSON{
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
