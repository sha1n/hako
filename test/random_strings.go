package test

import "math/rand"

const alphanumLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const alphanumLettersWithPunc = alphanumLetters + "!@#$%^&*()_-+={[}]|\\/?><,.~` "

// RandomStringN creates a new random alphanumeric string, 'n' characters long.
func RandomStringN(n int) string {
	return randomStringWith(alphanumLettersWithPunc, n)
}

// RandomString ...
func RandomString() string {
	return randomStringWith(alphanumLettersWithPunc, rand.Intn(50))
}

func randomStringWith(charRange string, length int) string {
	var letter = []rune(charRange)
	buffer := make([]rune, length)

	for i := range buffer {
		buffer[i] = letter[rand.Intn(len(letter))]
	}
	return string(buffer)
}
