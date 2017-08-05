package inputs

import (
	"errors"
	"fmt"
)

// validateFile validates an value against a file parameter
func (this _validator) validateFile(
	value *string,
) []error {
	if nil == value {
		return []error{errors.New("file required")}
	}

	fileInfo, err := this.os.Stat(*value)
	if nil != err {
		return []error{err}
	} else if !fileInfo.Mode().IsRegular() {
		return []error{fmt.Errorf("%v not a file", *value)}
	}
	return []error{}
}
