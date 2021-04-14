package helpers

import (
	"strings"
)

// StrSliceIndexOf returns index of given string in a slice or -1 if not found
func StrSliceIndexOf(slice []string, s string) int {
	for i, count := 0, len(slice); i < count; i++ {
		if slice[i] == s {
			return i
		}
	}
	return -1
}

func TrimWhiteSpaces(s string) string {

	return strings.Join(strings.Fields(s), " ")
}
