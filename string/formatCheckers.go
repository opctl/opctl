package string

import (
	"github.com/blang/semver"
	"github.com/docker/distribution/reference"
)

// DockerImageRefFormatChecker is a gojsonschema.FormatChecker for docker image refs
type DockerImageRefFormatChecker struct{}

// Implement gojsonschema.FormatChecker interface
func (f DockerImageRefFormatChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if ok == false {
		return false
	}

	_, err := reference.Parse(asString)
	return nil == err
}

// SemVerFormatChecker is a gojsonschema.FormatChecker for semantic versions
type SemVerFormatChecker struct{}

// Implement gojsonschema.FormatChecker interface
func (f SemVerFormatChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)
	if ok == false {
		return false
	}

	_, err := semver.Parse(asString)
	return nil == err
}
