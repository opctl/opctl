package data

import (
	"errors"
	"fmt"
)

// validateDir validates an value against a dir parameter
func (this _validator) validateDir(
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
