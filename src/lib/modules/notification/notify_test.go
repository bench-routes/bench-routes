package notify

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestSendNotification(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
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
	n.NotifyBasic()
	n.NotifyWithTimeSpan()
	n.NotifyWithImageIcon()
	n.NotifyWithUrgencyLevel()
	n.NotifyWithImageAndTimeSpan()
	n.NotifyWithUrgencyLevelAndTimeSpan()
}
