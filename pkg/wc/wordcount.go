package wc

import (
	"strings"
	"unicode"
)

// Count returns case-insensitive word frequencies for the provided text.
func Count(text string) map[string]int {
	counts := make(map[string]int)
	for _, word := range strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	}) {
		counts[strings.ToLower(word)]++
	}
	return counts
}
