package utils_test

import (
	"testing"

	utils "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
)

func TestHashPassword(t *testing.T) {
	t.Parallel()

	password := "testPassword"
	hash := utils.HashPassword(password)

	if hash == "" {
		t.Fatal("Hash should not be empty")
	}

	differentHash := utils.HashPassword("anotherPassword")
	if hash == differentHash {
		t.Fatal("Hashes for different passwords should be different")
	}
}

func TestCheckPassword(t *testing.T) {
	t.Parallel()

	password := "testPassword"
	hash := utils.HashPassword(password)

	if !utils.CheckPassword(password, hash) {
		t.Fatal("Password check should return true for correct password")
	}

	if utils.CheckPassword("wrongPassword", hash) {
		t.Fatal("Password check should return false for incorrect password")
	}
}
