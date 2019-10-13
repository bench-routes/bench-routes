package utils

import (
	"log"
	"net/http"
)

// VerifyConnection verifies the connectivity of the program with the external network
func VerifyConnection() (bool, int) {
	res, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		log.Printf("unable to ping external network\n")
		return false, -1
	}
	log.Printf("external connection verified\n")
	return true, res.StatusCode
}
