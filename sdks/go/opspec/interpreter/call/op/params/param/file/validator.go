package file

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/sdks/go/model"
)

type Validator interface {
	Validate(
		value *model.Value,
	) []error
}

func NewValidator() Validator {
	return _validator{
		os: ios.New(),
	}
}

type _validator struct {
	os ios.IOS
}

// Validate validates a value against a string parameter
func (vdt _validator) Validate(
	value *model.Value,
) []error {
	if nil == value || nil == value.File {
		return []error{errors.New("file required")}
	}

	fileInfo, err := vdt.os.Stat(*value.File)
	if nil != err {
		return []error{err}
	} else if !fileInfo.Mode().IsRegular() {
		return []error{fmt.Errorf("%v not a file", *value.File)}
	}
	return []error{}
}
