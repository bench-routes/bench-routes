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

// HTTPPingFilterValue filters the illegal characters by value,
// that may panic the ping running from the terminal.
func HTTPPingFilterValue(s string) string {
	tmp := s
	HTTPPingFilter(&tmp)
	s = tmp
	return s
}

// RouteDestroyer causes mayhem
func RouteDestroyer(url string) string {
	url = strings.ReplaceAll(url, "/", "_")
	url = strings.ReplaceAll(url, ":", "_")
	return url
}
