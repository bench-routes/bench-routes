package utils

import (
	"net/http"

	"github.com/zairza-cetb/bench-routes/src/lib/utils/logger"
)

// VerifyConnection verifies the connectivity of the program with the external network
func VerifyConnection() (bool, int) {
	res, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		logger.TerminalandFileLogger.Printf("unable to ping external network\n")
		return false, -1
	}
	logger.TerminalandFileLogger.Printf("external connection verified\n")
	return true, res.StatusCode
}
