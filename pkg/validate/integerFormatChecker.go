package validate

import (
	"strconv"
)

// Define the format checker
type IntegerFormatChecker struct{}

// Ensure it meets the gojsonschema.FormatChecker interface
func (f IntegerFormatChecker) IsFormat(input string) bool {
	_, err := strconv.Atoi(input)
	return nil == err
}
