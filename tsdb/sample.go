package main

import (
	"fmt"
	"time"
)

func main() {
	arr := []Block{}
	for i:=0; i < 10; i++ {
		inst := Block{}
		inst.Datapoint = i*100
		inst.NormalizedTime = time.Now().Unix()
		inst.Timestamp = time.Now()
		arr = append(arr, inst)
		if i == 0 {
			arr[i].PrevBlock = nil
		} else {
			arr[i].PrevBlock = &arr[i-1]
		}
		if i == 9 {
			arr[i].NextBlock = nil
		} else {
			arr[i-1].NextBlock = &arr[i]
		}
	}
	fmt.Println(arr)
}