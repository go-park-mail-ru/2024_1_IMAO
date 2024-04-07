package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func RandString(length int) string {
	var result string

	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result += string(letters[randomIndex.Int64()])
	}

	return result
}
