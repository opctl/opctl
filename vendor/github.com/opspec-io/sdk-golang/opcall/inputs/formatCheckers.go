package inputs

import (
	"github.com/blang/semver"
	"github.com/docker/distribution/reference"
	"strconv"
)

// DockerImageRefFormatChecker is a gojsonschema.FormatChecker for docker image refs
type DockerImageRefFormatChecker struct{}

// Implement gojsonschema.FormatChecker interface
func (f DockerImageRefFormatChecker) IsFormat(input string) bool {
	_, err := reference.Parse(input)
	return nil == err
}

// IntegerFormatChecker is a gojsonschema.FormatChecker for integers
type IntegerFormatChecker struct{}

// Implement gojsonschema.FormatChecker interface
func (f IntegerFormatChecker) IsFormat(input string) bool {
	_, err := strconv.Atoi(input)
	return nil == err
}

// SemVerFormatChecker is a gojsonschema.FormatChecker for semantic versions
type SemVerFormatChecker struct{}

// Implement gojsonschema.FormatChecker interface
func (f SemVerFormatChecker) IsFormat(input string) bool {
	_, err := semver.Parse(input)
	return nil == err
}
