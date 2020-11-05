package dir

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
	if nil == value || nil == value.Dir {
		return []error{errors.New("dir required")}
	}

	fileInfo, err := os.Stat(*value.Dir)
	if nil != err {
		return []error{err}
	} else if !fileInfo.IsDir() {
		return []error{fmt.Errorf("%v not a dir", *value.Dir)}
	}
	return []error{}
}
