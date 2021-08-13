package opgraph

import (
	"regexp"
	"unicode/utf8"

	"github.com/fatih/color"
)

var muted = color.New(color.Faint)
var highlighted = color.New(color.Bold)
var success = color.New(color.FgGreen)
var failed = color.New(color.FgRed)
var warning = color.New(color.FgYellow)

var ansiRegex = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

func stripAnsi(str string) string {
	return ansiRegex.ReplaceAllString(str, "")
}

func stripAnsiToLength(str string, length int) string {
	for countChars(stripAnsi(str)) > length {
		_, size := utf8.DecodeLastRuneInString(str)
		str = str[:len(str)-size]
	}
	return str
}

func countChars(str string) int {
	count := 0
	for len(str) > 0 {
		_, size := utf8.DecodeLastRuneInString(str)
		count++
		str = str[:len(str)-size]
	}
	return count
}
