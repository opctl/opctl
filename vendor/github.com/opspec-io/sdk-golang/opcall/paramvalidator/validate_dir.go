package paramvalidator

import (
	"errors"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
)

// validateDir validates an value against a dir parameter
func (this _ParamValidator) validateDir(
	rawValue *string,
	param *model.DirParam,
) []error {

	value := rawValue
	if nil == value && nil != param.Default {
		// apply default if value not set
		value = param.Default
	}

	if nil == value {
		return []error{errors.New("Dir required")}
	}

	fileInfo, err := this.os.Stat(*value)
	if nil != err {
		return []error{err}
	} else if !fileInfo.IsDir() {
		return []error{fmt.Errorf("%v not a dir", *value)}
	}
	return []error{}
}
