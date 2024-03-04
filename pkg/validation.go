package pkg

import (
	"net/mail"
	"strings"
	"unicode"
)

const (
	ErrTooShort    = "Password is too short"
	ErrTooLong     = "Password is too long"
	ErrWrongFormat = "Password does not contain specific symbols"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func ValidatePassword(password string) string {
	switch {
	case len(password) < 8:
		return ErrTooShort
	case len(password) > 32:
		return ErrTooLong
	}

	specialChars := "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	if !(hasLower && hasUpper && hasSpecial && hasDigit) {
		return ErrWrongFormat
	}

	return ""
}
