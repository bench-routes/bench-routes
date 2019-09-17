package scraps

import (
	"testing"
)

func TestCLIJitterScrap(t *testing.T) {
	for _, samples := range input {
		CLIJitterScrap(&samples)
	}
}