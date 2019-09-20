package utils

import (
	"testing"
)

var (
	// to verify whether the service is able to respond and request to multiple website urls
	urlsDiff = []string{
		"https://www.google.co.in/",
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
		chnlCLICommunication := make(chan *string)
		go CLIPing(&ele, 2, chnlCLICommunication)
		resp := *(<-chnlCLICommunication)
		if len(resp) == 0 {
			t.Errorf("err requesting %s\n", ele)
		} else {
			t.Logf("%s\n", resp)
		}
	}

	// testing packets on permutative urls
	for _, ele := range urlsPermute {
		chnlCLICommunication := make(chan *string)
		go CLIPing(&ele, 2, chnlCLICommunication)
		resp := *(<-chnlCLICommunication)
		if len(resp) == 0 {
			t.Errorf("err requesting %s\n", ele)
		}
	}
}
