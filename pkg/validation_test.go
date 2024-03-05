package pkg_test

import (
	"testing"

	"github.com/go-park-mail-ru/2024_1_IMAO/pkg"
)

func TestValidateEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"Valid email", "test@example.com", true},
		{"Invalid email", "test@.com", false},
		{"Invalid email", "test.com", false},
		{"Invalid email", "test@", false},
		{"Empty email", "", false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := pkg.ValidateEmail(tt.email); got != tt.want {
				t.Errorf("ValidateEmail(%s) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestCheckSymbols(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "Valid password",
			password: "Password123!",
			want:     true,
		},
		{
			name:     "Invalid password - missing special character",
			password: "Password123",
			want:     false,
		},
		{
			name:     "Invalid password - missing uppercase",
			password: "password123!",
			want:     false,
		},
		{
			name:     "Invalid password - missing lowercase",
			password: "PASSWORD123!",
			want:     false,
		},
		{
			name:     "Invalid password - missing digit",
			password: "Password!",
			want:     false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := pkg.CheckSymbols(tt.password); got != tt.want {
				t.Errorf("CheckSymbols() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		password string
		want     string
	}{
		{
			name:     "Valid password",
			password: "Password123!",
			want:     "",
		},
		{
			name:     "Too short password",
			password: "Pass!",
			want:     pkg.ErrTooShort,
		},
		{
			name:     "Too long password",
			password: "ThisPasswordIsWayTooLongAndDoesNotMeetTheLengthRequirements123!",
			want:     pkg.ErrTooLong,
		},
		{
			name:     "Invalid format password",
			password: "password",
			want:     pkg.ErrWrongFormat,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := pkg.ValidatePassword(tt.password); got != tt.want {
				t.Errorf("ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
