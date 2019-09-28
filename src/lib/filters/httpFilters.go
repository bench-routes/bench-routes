package filters

import (
	"strings"
)

// HTTPPingFilter filters the illegal characters that may panic the ping
// subprocess running from the terminal
func HTTPPingFilter(s *string) *string {
	urlParts := strings.Split(*s, "/")
	for _, urlPart := range urlParts {
		if strings.Contains(urlPart, ".") {
			*s = urlPart
		}
	}

	*s = strings.Replace(*s, "www.", "", -1)
	*s = strings.Replace(*s, "https:", "", -1)
	*s = strings.Replace(*s, "http:", "", -1)
	*s = strings.Split(*s, ":")[0]

	return s
}
