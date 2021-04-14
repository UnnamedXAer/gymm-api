package helpers

import (
	"strings"
	"unicode"
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
func TrimWhiteSpacesBuilder(s string) string {
	panic("not implemented yet.")
	var builder strings.Builder
	builder.Grow(len(s))
	var lastRune rune
	for _, r := range s {
		if unicode.IsSpace(r) {
			if unicode.IsSpace(lastRune) {
				continue
			}
			lastRune = ' '
			builder.WriteRune(' ')
			continue
		}
		lastRune = r
		builder.WriteRune(r)
	}

	return builder.String()
}
