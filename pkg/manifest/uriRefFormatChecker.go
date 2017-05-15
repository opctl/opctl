package manifest

import (
	"net/url"
)

// Define the format checker
type uriRefFormatChecker struct{}

// Ensure it meets the gojsonschema.FormatChecker interface
func (f uriRefFormatChecker) IsFormat(input string) bool {
	//@TODO delete once native support exists. 'uri-reference' is a json schema v6 keyword; gojsonschema currently supports v4
	if _, err := url.Parse(input); nil != err {
		return false
	}

	return true
}
