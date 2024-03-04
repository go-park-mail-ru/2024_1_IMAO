package pkg

import "math/rand"

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func RandString(length int) string {
	var result string
	for i := 0; i < length; i++ {
		result += string(letters[rand.Intn(len(letters))])
	}

	return result
}
