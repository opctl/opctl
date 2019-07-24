package param

import (
	"errors"
	"github.com/opctl/opctl/sdks/go/types"
)

// validateBoolean validates a value against a boolean parameter
func (vdt _validator) validateBoolean(
	value *types.Value,
) []error {

	// handle no value passed
	if nil == value || nil == value.Boolean {
		return []error{errors.New("boolean required")}
	}

	return []error{}
}
