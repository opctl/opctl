package validate

import (
	"errors"
	"github.com/opspec-io/sdk-golang/model"
)

// validates an value against a network socket parameter
func (this validate) socketParam(
	rawValue *model.Data,
	param *model.SocketParam,
) (errs []error) {
	errs = []error{}

	// handle no value passed
	if nil == rawValue || "" == rawValue.Socket {
		errs = append(errs, errors.New("Socket required"))
	}
	return
}
