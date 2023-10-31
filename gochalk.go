package gochalk

import (
	"fmt"
	"strings"
)

const escape = "\x1b"
const resetColor = escape + "[0m"

const (
	FgBlack int = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

const (
	FgBrightBlack int = iota + 90
	FgBrightRed
	FgBrightGreen
	FgBrightYellow
	FgBrightBlue
	FgBrightMagenta
	FgBrightCyan
	FgBrightWhite
)

const (
	BgBlack int = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

const (
	BgBrightBlack int = iota + 100
	BgBrightRed
	BgBrightGreen
	BgBrightYellow
	BgBrightBlue
	BgBrightMagenta
	BgBrightCyan
	BgBrightWhite
)

// Method to print string with red foreground color
func Red(value ...string) string {
	return getSingleString(FgRed, value...)
}

// Method to print string with green foreground color
func Green(value ...string) string {
	return getSingleString(FgGreen, value...)
}

// Method to print string with yellow foreground color
func Yellow(value ...string) string {
	return getSingleString(FgYellow, value...)
}

// Method to print string with blue foreground color
func Blue(value ...string) string {
	return getSingleString(FgBlue, value...)
}

// Method to print string with magenta foreground color
func Magenta(value ...string) string {
	return getSingleString(FgMagenta, value...)
}

// Method to print string with cyan foreground color
func Cyan(value ...string) string {
	return getSingleString(FgCyan, value...)
}

// Method to print string with white foreground color
func White(value ...string) string {
	return getSingleString(FgWhite, value...)
}

func getSingleString(color int, strs ...string) string {
	if len(strs) == 0 {
		return ""
	}

	var finalString string
	for index, str := range strs {
		if strings.Contains(str, escape) && index != len(strs)-1 {
			finalString += str + escapedColor(color) + " "
		} else if index == len(strs)-1 {
			finalString += str
		} else {
			finalString += str + " "
		}
	}

	return getColoredString(color, finalString)
}

// Method to return color in escaped string format
func escapedColor(color int) string {
	return fmt.Sprintf("%s[%dm", escape, color)
}

// Method to return given string encapsulated with color styling
func getColoredString(color int, value string) string {
	return fmt.Sprintf("%s%s%s", escapedColor(color), value, resetColor)
}
