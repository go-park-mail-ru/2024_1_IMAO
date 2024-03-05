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

const (
	minLen = 8
	maxLen = 32
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func checkSymbols(password string) bool {
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

	return hasLower && hasUpper && hasSpecial && hasDigit
}

func ValidatePassword(password string) string {
	switch {
	case len(password) < minLen:
		return ErrTooShort
	case len(password) > maxLen:
		return ErrTooLong
	}

	if !checkSymbols(password) {
		return ErrWrongFormat
	}

	return ""
}
