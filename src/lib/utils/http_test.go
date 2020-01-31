package utils

import (
	"testing"
)

var (
	// to verify whether the service is able to respond and request to multiple website urls
	urlsDiff = []string{
		"https://www.google.co.in/",
		"https://www.google.co.in/search?q=hello",
		"https://www.facebook.com/",
		"https://www.yahoo.com/",
	}
	// to verify whether the filters are working as expected or not
	urlsPermute = []string{
		"https:www.google.co.in",
		"http:www.google.co.in",
		"google.co.in/",
		"//www.google.co.in",
	}
)

func TestCLIPing(t *testing.T) {
	// testing packets on diff urls
	for _, ele := range urlsDiff {
		_, err := CLIPing(ele, 2)
		if err != nil {
			t.Errorf("err requesting %s\n", ele)
		}
	}

	// testing packets on permutative urls
	for _, ele := range urlsPermute {
		_, err := CLIPing(ele, 2)
		if err != nil {
			t.Errorf("err requesting %s\n", ele)
		}
	}
}
