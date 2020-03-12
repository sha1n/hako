package utils

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
)

func Test_RandomString(t *testing.T) {
	length := randomLength()
	randomString := RandomString(length)

	assert.Len(t, randomString, length)
	assertInRange(t, randomString, alphanumLettersWithPunc)
}

func Test_RandomStr50(t *testing.T) {
	assert.Len(t, RandomStr50(), 50)
	assertInRange(t, RandomStr50(), alphanumLettersWithPunc)
}

func Test_RandomAlphaNumericString(t *testing.T) {
	length := randomLength()
	randomString := RandomAlphaNumericString(length)
	assert.Len(t, RandomAlphaNumericString(length), length)
	assertInRange(t, randomString, alphanumLetters)
}

func assertInRange(t *testing.T, randomString string, charRange string) {
	for _, runeValue := range randomString {
		assert.True(t, strings.Contains(charRange, string(runeValue)))
	}
}

func randomLength() int {
	return rand.Intn(100)
}
