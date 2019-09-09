package utils

import (
	"testing"
)

func TestVerifyConnection(t *testing.T) {
	runner, code := VerifyConnection()
	if !runner {
		t.Errorf("connection to external network failed\n")
	} else {
		t.Logf("%d status code from verification of external network\n", code)
	}
}
