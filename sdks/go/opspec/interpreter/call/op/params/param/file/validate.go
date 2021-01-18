package file

import (
	"errors"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"os"
)

// Validate validates a value against a string parameter
func Validate(
	value *model.Value,
) []error {
	if nil == value || nil == value.File {
		return []error{errors.New("file required")}
	}

	fileInfo, err := os.Stat(*value.File)
	if nil != err {
		return []error{err}
	} else if !fileInfo.Mode().IsRegular() {
		return []error{fmt.Errorf("%v not a file", *value.File)}
	}
	return []error{}
}
