package validator

import (
	"errors"
	"github.com/opspec-io/sdk-golang/model"
)

// validateSocket validates an value against a socket parameter
func (this _Validator) validateSocket(
	rawValue *string,
	param *model.SocketParam,
) (errs []error) {
	errs = []error{}

	// handle no value passed
	if nil == rawValue {
		errs = append(errs, errors.New("Socket required"))
	}
	return
}
