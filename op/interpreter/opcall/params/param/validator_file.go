package param

import (
	"errors"
	"fmt"
	"github.com/opctl/sdk-golang/model"
)

// validateFile validates a value against a file parameter
func (vdt _validator) validateFile(
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
