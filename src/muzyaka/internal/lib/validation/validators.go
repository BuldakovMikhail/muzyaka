package validation

import (
	"unicode"
)

func ValidateWithoutSpace(str string) bool {
	for _, c := range str {
		if unicode.IsSpace(c) {
			return false
		}
	}
	return true
}
