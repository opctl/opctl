package dir

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ Validator

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
)

type Validator interface {
	// Validate validates a dir
	Validate(
		value *string,
	) []error
}

func newValidator() Validator {
	return _validator{
		os: ios.New(),
	}
}

type _validator struct {
	os ios.IOS
}

func (this _validator) Validate(
	value *string,
) []error {
	if nil == value {
		return []error{errors.New("dir required")}
	}

	fileInfo, err := this.os.Stat(*value)
	if nil != err {
		return []error{err}
	} else if !fileInfo.IsDir() {
		return []error{fmt.Errorf("%v not a dir", *value)}
	}
	return []error{}
}
