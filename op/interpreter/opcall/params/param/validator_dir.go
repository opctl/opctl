package param

import (
	"errors"
	"fmt"
	"github.com/opctl/sdk-golang/model"
)

// validateDir validates a value against a dir parameter
func (vdt _validator) validateDir(
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
