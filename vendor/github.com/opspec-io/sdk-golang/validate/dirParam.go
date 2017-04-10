package validate

import (
	"errors"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
)

// validates an value against a dir parameter
func (this validate) dirParam(
	rawValue *model.Data,
	param *model.DirParam,
) []error {

	// handle no value passed
	if nil == rawValue || "" == rawValue.Dir {
		return []error{errors.New("Dir required")}
	}

	fileInfo, err := this.fs.Stat(rawValue.Dir)
	if nil != err {
		return []error{err}
	} else if !fileInfo.IsDir() {
		return []error{fmt.Errorf("%v not a dir", rawValue.Dir)}
	}
	return []error{}
}
