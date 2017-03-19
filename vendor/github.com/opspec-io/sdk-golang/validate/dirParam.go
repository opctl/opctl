package validate

import (
	"errors"
	"github.com/opspec-io/sdk-golang/model"
)

// validates an value against a dir parameter
func (this validate) dirParam(
	rawValue *model.Data,
	param *model.DirParam,
) (errs []error) {
	errs = []error{}

	// handle no value passed
	if nil == rawValue || "" == rawValue.Dir {
		errs = append(errs, errors.New("Dir required"))
	}
	return
}
