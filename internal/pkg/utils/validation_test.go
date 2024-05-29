package utils_test

import (
	"errors"
	"testing"

	utils "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
)

var (
	ErrorOne = errors.New("mail: no angle-addr")
)

func TestValidateEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		email string
		err   error
	}{
		{"Valid email", "test@example.com", nil},
		{"Invalid email", "test@.com", ErrorOne},
		{"Invalid email", "test.com", errors.New("mail: missing '@' or angle-addr")},
		{"Invalid email", "test@", errors.New("mail: no angle-addr")},
		{"Empty email", "", errors.New("mail: no address")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.name == "Valid email" {

				if got := utils.ValidateEmail(tt.email); got != nil {

					t.Errorf("ValidateEmail(%s) = %v, want %v", tt.email, got, tt.err)
				}

			} else {

				if got := utils.ValidateEmail(tt.email); got == nil {

					t.Errorf("ValidateEmail(%s) = %v, want %v", tt.email, got, tt.err)
				}
			}

		})
	}
}

func compareSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		password string
		email    string
		want     []string
	}{
		{
			name:     "Valid password",
			password: "Password123!",
			email:    "test@example.com",
			want:     []string{},
		},
		{
			name:     "Too short password",
			password: "Pass!",
			email:    "test@.com",
			want:     []string{"Password is too short", "Password does not contain required symbols", "Wrong email format"},
		},
		{
			name:     "Too long password",
			password: "ThisPasswordIsWayTooLongAndDoesNotMeetTheLengthRequirements123!",
			email:    "test.com",
			want:     []string{"Password is too long", "Wrong email format"},
		},
		{
			name:     "Invalid format password",
			password: "password",
			email:    "test@",
			want:     []string{"Password does not contain required symbols", "Wrong email format"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := utils.Validate(tt.email, tt.password); !compareSlices(got, tt.want) {
				t.Errorf("ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
