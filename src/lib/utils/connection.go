package utils

import (
	"net/http"

	"github.com/zairza-cetb/bench-routes/src/lib/utils/logger"
)

// VerifyConnection verifies the connectivity of the program with the external network
func VerifyConnection() (bool, int) {
	res, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		logger.Terminal("unable to ping external network", "p")
		return false, -1
	}
	logger.File("external connection verified", "p")
	return true, res.StatusCode
}
