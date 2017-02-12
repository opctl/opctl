package validate

import (
	"github.com/docker/distribution/reference"
)

// Define the format checker
type DockerImageRefFormatChecker struct{}

// Ensure it meets the gojsonschema.FormatChecker interface
func (f DockerImageRefFormatChecker) IsFormat(input string) bool {
	_, err := reference.Parse(input)
	return nil == err
}
