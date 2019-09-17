package jitter

import (
	"testing"
)

var (
	urls = []string{
		"google.co.in",
		"facebook.com",
		"yahoo.com",
		"github.com",
		"youtube.com",
	}
)

func TestHandlerJitter(t *testing.T) {
	for _, inst := range urls {
		HandleJitter(&inst, 20)
	}
}