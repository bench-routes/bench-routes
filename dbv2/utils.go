package dbv2

import "fmt"

func ConvertValueToValueSet(data ...string) string {
	var s string
	for i := range data {
		s = fmt.Sprintf("%s%s%c", s, data[i], valueSeparator)
	}
	return s[:len(s)-1] // Ignore the last pipe.
}
