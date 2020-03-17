package console

import "fmt"

const RESET = "\u001b[0m"
const RED = "\u001b[31m"
const GREEN = "\u001b[32m"
const YELLOW = "\u001b[33m"
const CYAN = "\u001b[36m"
const BOLD = "\u001b[1m"
const REVERSE = "\u001b[7m"

func Reverse(s string, args ...interface{}) string {
	return formatAndReset(fmt.Sprintf(s, args...), REVERSE)
}

func Bold(s string, args ...interface{}) string {
	return formatAndReset(fmt.Sprintf(s, args...), BOLD)
}

func Green(s string, args ...interface{}) string {
	return formatAndReset(fmt.Sprintf(s, args...), GREEN)
}

func Yellow(s string, args ...interface{}) string {
	return formatAndReset(fmt.Sprintf(s, args...), YELLOW)
}

func Red(s string, args ...interface{}) string {
	return formatAndReset(fmt.Sprintf(s, args...), RED)
}

func Cyan(s string, args ...interface{}) string {
	return formatAndReset(fmt.Sprintf(s, args...), CYAN)
}

func formatAndReset(s string, styleCode string) string {
	return fmt.Sprintf("%s%s%s", styleCode, s, RESET)
}
