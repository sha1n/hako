package console

import "fmt"

const reset = "\u001b[0m"
const red = "\u001b[31m"
const green = "\u001b[32m"
const yellow = "\u001b[33m"
const cyan = "\u001b[36m"
const bold = "\u001b[1m"
const reserve = "\u001b[7m"

// Reverse formats the input string and sets its background color to its foreground color and vice versa.
func Reverse(s string, args ...interface{}) string {
	return formatAndReset(fmt.Sprintf(s, args...), reserve)
}

// Bold formats the input string and sets its style to bold.
func Bold(s string, args ...interface{}) string {
	return formatAndReset(fmt.Sprintf(s, args...), bold)
}

// Green formats the input string and sets its color to green.
func Green(s string, args ...interface{}) string {
	return formatAndReset(fmt.Sprintf(s, args...), green)
}

// Yellow formats the input string and sets its color to yellow.
func Yellow(s string, args ...interface{}) string {
	return formatAndReset(fmt.Sprintf(s, args...), yellow)
}

// Red formats the input string and sets its color to red.
func Red(s string, args ...interface{}) string {
	return formatAndReset(fmt.Sprintf(s, args...), red)
}

// Cyan formats the input string and sets its color to cyan.
func Cyan(s string, args ...interface{}) string {
	return formatAndReset(fmt.Sprintf(s, args...), cyan)
}

func formatAndReset(s string, styleCode string) string {
	return fmt.Sprintf("%s%s%s", styleCode, s, reset)
}
