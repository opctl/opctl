package validate

import (
	"errors"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
)

// validates an value against a file parameter
func (this validate) fileParam(
	rawValue *model.Data,
	param *model.FileParam,
) []error {
	// handle no value passed
	if nil == rawValue || "" == rawValue.File {
		return []error{errors.New("File required")}
	}

	fileInfo, err := this.fs.Stat(rawValue.File)
	if nil != err {
		return []error{err}
	} else if !fileInfo.Mode().IsRegular() {
		return []error{fmt.Errorf("%v not a file", rawValue.File)}
	}
	return []error{}
}
