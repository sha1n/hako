package console

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var styleFunctions = [...]func(string, ...interface{}) string{
	Bold,
	Reverse,
	Green,
	Yellow,
	Red,
	Cyan,
}

var allAsciiCodes = reset + reserve + bold + green + yellow + red + cyan

func Test_FormatArgumentsPassThrough(t *testing.T) {

	for mi := range styleFunctions {
		formattedString := styleFunctions[mi]("Hello %s", "World")
		strippedString := strings.Trim(formattedString, allAsciiCodes)
		assert.Equal(t, "Hello World", strippedString)
	}
}

func Test_AllEndWithResetCode(t *testing.T) {

	for mi := range styleFunctions {
		formattedString := styleFunctions[mi]("Hello %s", "World")
		assert.True(t, strings.HasSuffix(formattedString, reset))
	}
}
