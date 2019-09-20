package filters

import (
	"strings"
)

// HTTPPingFilter filters the illegal characters that may panic the ping
// subprocess running from the terminal
func HTTPPingFilter(s *string) *string {
	*s = strings.Replace(*s, "http://", "", -1)
	*s = strings.Replace(*s, "https://", "", -1)
	*s = strings.Replace(*s, "https:", "", -1)
	*s = strings.Replace(*s, "http:", "", -1)
	*s = strings.Replace(*s, "/", "", -1)
	*s = strings.Replace(*s, "www.", "", -1)
	*s = strings.Replace(*s, ":", "", -1)
	do := [...]string{".com", ".in", ".org", ".co.in", ".edu"}
	
	temp := *s
	for _, value := range do {
		if (strings.Index(*s, value) > 0) {
			v := (strings.Index(*s, value))
			temp = temp[:v+len(value)]
			*s = temp
		}
		
	}
	return s
}
