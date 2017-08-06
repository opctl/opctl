package inputs

import (
	"errors"
)

// validateSocket validates an value against a socket parameter
func (this _validator) validateSocket(
	rawValue *string,
) []error {

	// handle no value passed
	if nil == rawValue {
		return []error{errors.New("socket required")}
	}

	return []error{}
}
