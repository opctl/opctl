package validate

import (
	"errors"
	"github.com/opspec-io/sdk-golang/model"
)

// validates an value against a file parameter
func (this validate) fileParam(
	rawValue *model.Data,
	param *model.FileParam,
) (errs []error) {
	errs = []error{}

	// handle no value passed
	if nil == rawValue || "" == rawValue.File {
		errs = append(errs, errors.New("File required"))
	}
	return
}
