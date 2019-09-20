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
	return s
}
