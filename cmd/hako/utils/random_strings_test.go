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

func assertInRange(t *testing.T, randomString string, charRange string) {
	for _, runeValue := range randomString {
		assert.True(t, strings.Contains(charRange, string(runeValue)))
	}
}

func randomLength() int {
	return rand.Intn(100)
}
