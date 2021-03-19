package file

import (
	"errors"
	"fmt"
	"os"

	"github.com/opctl/opctl/sdks/go/model"
)

// Validate validates a value against a string parameter
func Validate(
	value *model.Value,
) []error {
	if value == nil || value.File == nil {
		return []error{errors.New("file required")}
	}

	fileInfo, err := os.Stat(*value.File)
	if err != nil {
		return []error{err}
	} else if !fileInfo.Mode().IsRegular() {
		return []error{fmt.Errorf("%v not a file", *value.File)}
	}
	return []error{}
}
