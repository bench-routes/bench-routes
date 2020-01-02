package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	file, ok := SetupLogger()
	if !ok {
		t.Errorf("Setting up the logger failed")
	} else {
		t.Logf("Logger file created")
		t.Logf("%s", file.Name())
	}
}
