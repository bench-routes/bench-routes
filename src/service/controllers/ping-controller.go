package controllers

import (
	"fmt"
)

// PingController controllers the ping requests and transfers to the respective handler
func PingController() bool {
	fmt.Println("ping controller")
	return true
}
