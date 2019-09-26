package notify

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func TestNotify(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	substr := "/bench-routes/"
	c := strings.Index(dir, substr)
	dir = dir[:c+len(substr)] + "assets/icon.png"
	var n = NotificationData{
		title:       "hi",
		description: "hello",
		icon:        dir,
		urgency:     "critical",
		time:        5000,
	}
	n.Notify1()
	n.Notify2()
	n.Notify3()
	n.Notify4()
	n.Notify5()
	n.Notify6()
}
