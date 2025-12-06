package utils

import (
	"strings"
)

func ValidPassword(password string) bool {
	specialChars := "!@#$%()"

	for _, i := range password {
		if strings.ContainsRune(specialChars, i) {
			return true
		}
	}

	return false
}

func ValidLen(password string) bool {
	if len(password) > 8 {
		return true
	} else {
	return false
	}
}

// func IsEnglishOnly(password string) bool {
// 	for _, r := range password {
// 		if !unicode.IsLetter(r) || (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
// 			return false 
// 		}
// 	}
// 	return true
// }
