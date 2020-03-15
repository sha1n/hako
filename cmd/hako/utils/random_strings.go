package utils

import "math/rand"

const alphanumLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const alphanumLettersWithPunc = alphanumLetters + "!@#$%^&*()_-+={[}]|\\/?><,.~` "

func RandomString(n int) string {
	return randomAlphNumericString(alphanumLettersWithPunc, n)
}

func RandomStr50() string {
	return RandomString(50)
}

func RandomAlphaNumericString(n int) string {
	return randomAlphNumericString(alphanumLetters, n)
}

func randomAlphNumericString(charRange string, n int) string {
	var letter = []rune(charRange)
	buffer := make([]rune, n)

	for i := range buffer {
		buffer[i] = letter[rand.Intn(len(letter))]
	}
	return string(buffer)
}
