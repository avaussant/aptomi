package util

import (
	"strings"
)

// EscapeName escapes provided string by replacing # and _ with -
func EscapeName(str string) string {
	r := strings.NewReplacer("#", "-", "_", "-")
	return r.Replace(str)
}

// StringContainsAny returns if string contains any of the given substrings
func StringContainsAny(str string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(str, substr) {
			return true
		}
	}

	return false
}
