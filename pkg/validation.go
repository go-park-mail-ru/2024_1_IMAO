package pkg

import (
	"net/mail"
	"unicode"
)

const (
	ErrTooShort         = "Password is too short"
	ErrTooLong          = "Password is too long"
	ErrWrongFormat      = "Password does not contain required symbols"
	ErrWrongEmailFormat = "Wrong email format"
)

const (
	minLen = 8
	maxLen = 32
)

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func checkSymbols(password string) bool {
	var hasUpper, hasLower, hasDigit bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasLower && hasUpper && hasDigit
}

func validatePassword(password string) []string {
	var errors []string

	switch {
	case len(password) < minLen:
		errors = append(errors, ErrTooShort)
	case len(password) > maxLen:
		errors = append(errors, ErrTooLong)
	}

	if !checkSymbols(password) {
		errors = append(errors, ErrWrongFormat)
	}

	return errors
}

func Validate(email, password string) []string {
	errors := validatePassword(password)

	if !validateEmail(email) {
		errors = append(errors, ErrWrongEmailFormat)
	}

	return errors
}
