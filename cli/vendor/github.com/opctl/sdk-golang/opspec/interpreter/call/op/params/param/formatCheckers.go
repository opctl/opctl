package param

import (
	"github.com/blang/semver"
	"github.com/docker/distribution/reference"
	"math/big"
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

// IntegerFormatChecker is a gojsonschema.FormatChecker for integers
type IntegerFormatChecker struct{}

// Implement gojsonschema.FormatChecker interface
func (f IntegerFormatChecker) IsFormat(input interface{}) bool {

	asRat, ok := input.(*big.Rat)
	if ok == false {
		return false
	}

	return asRat.IsInt()
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
