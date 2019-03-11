package param

import (
	"errors"
	"github.com/opctl/sdk-golang/model"
)

// validateSocket validates a value against a socket parameter
func (vdt _validator) validateSocket(
	value *model.Value,
) []error {

	// handle no value passed
	if nil == value || nil == value.Socket {
		return []error{errors.New("socket required")}
	}

	return []error{}
}
