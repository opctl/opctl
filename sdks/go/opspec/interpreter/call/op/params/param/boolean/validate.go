package boolean

import (
	"errors"
	"github.com/opctl/opctl/sdks/go/model"
)

// Validate validates a value against a string parameter
func Validate(
	value *model.Value,
) []error {
	// handle no value passed
	if nil == value || nil == value.Boolean {
		return []error{errors.New("boolean required")}
	}

	return []error{}
}
