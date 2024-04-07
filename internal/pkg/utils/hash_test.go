package utils_test

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	t.Parallel()

	password := "testPassword"
	hash := HashPassword(password)

	if hash == "" {
		t.Fatal("Hash should not be empty")
	}

	differentHash := HashPassword("anotherPassword")
	if hash == differentHash {
		t.Fatal("Hashes for different passwords should be different")
	}
}

func TestCheckPassword(t *testing.T) {
	t.Parallel()

	password := "testPassword"
	hash := HashPassword(password)

	if !CheckPassword(password, hash) {
		t.Fatal("Password check should return true for correct password")
	}

	if CheckPassword("wrongPassword", hash) {
		t.Fatal("Password check should return false for incorrect password")
	}
}
