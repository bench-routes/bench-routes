package ping

import (
	"log"
	"reflect"
	"testing"

	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
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
		a := HandleFloodPing(&inst, 1000)
		if reflect.TypeOf(a).String() != "uint64" {
			t.Errorf("invalid response from HandleFloodPing")
		} else {
			log.Println(a)
		}
	}
}
