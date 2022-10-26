package isogram

import (
	"strings"
	"unicode"
)

func IsIsogram(word string) bool {
	m := make(map[string]int)
	for _, r := range word {
		if (unicode.IsLetter(r)) {
			c := string(r)
			c_lower := strings.ToLower(c)
			m[c_lower] += 1
		}
	}
	result := true
	for _, v := range m {
		if (v > 1) {
			result = false
		}
	}
	return result
}
