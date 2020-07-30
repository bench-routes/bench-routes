package utils

import (
	"net/http"

	"github.com/prometheus/common/log"
)

// VerifyConnection verifies the connectivity of the program with the external network
func VerifyConnection() (bool, int) {
	res, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		log.Infoln("unable to ping external network")
		return false, -1
	}
	return true, res.StatusCode
}
