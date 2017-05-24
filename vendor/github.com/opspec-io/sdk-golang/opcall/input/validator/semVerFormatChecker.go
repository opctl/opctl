package validator

import (
	"github.com/blang/semver"
)

// Define the format checker
type SemVerFormatChecker struct{}

// Ensure it meets the gojsonschema.FormatChecker interface
func (f SemVerFormatChecker) IsFormat(input string) bool {
	_, err := semver.Parse(input)
	return nil == err
}
