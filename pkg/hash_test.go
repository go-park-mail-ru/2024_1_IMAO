package pkg_test

import (
	"testing"

	"github.com/go-park-mail-ru/2024_1_IMAO/pkg"
)

func TestHashPassword(t *testing.T) {
	t.Parallel()

	password := "testPassword"
	hash := pkg.HashPassword(password)

	if hash == "" {
		t.Fatal("Hash should not be empty")
	}

	differentHash := pkg.HashPassword("anotherPassword")
	if hash == differentHash {
		t.Fatal("Hashes for different passwords should be different")
	}
}

func TestCheckPassword(t *testing.T) {
	t.Parallel()

	password := "testPassword"
	hash := pkg.HashPassword(password)

	if !pkg.CheckPassword(password, hash) {
		t.Fatal("Password check should return true for correct password")
	}

	if pkg.CheckPassword("wrongPassword", hash) {
		t.Fatal("Password check should return false for incorrect password")
	}
}
