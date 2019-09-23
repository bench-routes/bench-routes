package ping

import (
	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
	"log"
	"reflect"
	"testing"
)

var (
	urls = []string{
		"google.co.in",
		"facebook.com",
		"yahoo.com",
		"youtube.com",
	}
)

func TestHandlerPing(t *testing.T) {
	for _, inst := range urls {
		a := HandlePing(&inst, 20)
		if *a == (scrap.TypePingScrap{}) {
			t.Errorf("invalid response from HandlePing")
		} else {
			log.Println(*a)
		}
	}
}

func TestHandlerFloodPing(t *testing.T) {
	for _, inst := range urls {
		a, b := HandleFloodPing(&inst, 1000)
		if reflect.TypeOf(a).String() != "float64" {
			t.Errorf("invalid packet loss type, not a float64 value")
		} else if a > 100 || a < 0 {
			t.Errorf("invalid packet response")
		} else {
			log.Println(a)
		}
		if *b == (scrap.TypePingScrap{}) {
			t.Errorf("invalid ping response")
		} else {
			log.Println(*b)
		}
	}
}
