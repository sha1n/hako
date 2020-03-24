package utils

import "math/rand"

const alphanumLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const alphanumLettersWithPunc = alphanumLetters + "!@#$%^&*()_-+={[}]|\\/?><,.~` "

// RandomString creates a new random alphanumeric string, 'n' characters long.
func RandomString(n int) string {
	return randomStringWith(alphanumLettersWithPunc, n)
}

func randomStringWith(charRange string, length int) string {
	var letter = []rune(charRange)
	buffer := make([]rune, length)

	for i := range buffer {
		buffer[i] = letter[rand.Intn(len(letter))]
	}
	return string(buffer)
}
