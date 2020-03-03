package dir

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
	if nil == value || nil == value.Dir {
		return []error{errors.New("dir required")}
	}

	fileInfo, err := vdt.os.Stat(*value.Dir)
	if nil != err {
		return []error{err}
	} else if !fileInfo.IsDir() {
		return []error{fmt.Errorf("%v not a dir", *value.Dir)}
	}
	return []error{}
}
