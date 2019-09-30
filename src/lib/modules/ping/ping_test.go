package ping

import (
	"log"
	"reflect"
	"sync"
	"testing"

	scrap "github.com/zairza-cetb/bench-routes/src/lib/filters/scraps"
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
	var wg sync.WaitGroup
	wg.Add(len(urls))
	for _, inst := range urls {
		go HandlePing(&inst, 20, &wg)
	}
	wg.Wait()
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
