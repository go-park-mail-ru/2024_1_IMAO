package usecase

import (
	"net/mail"
	"strings"
	"unicode"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func ValidatePassword(password string) bool {
	if len(password) < 8 || len(password) > 32 {
		return false
	}

	specialChars := "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else if strings.ContainsRune(specialChars, char) {
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}
