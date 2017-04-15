package validate

import (
	"errors"
	"github.com/opspec-io/sdk-golang/model"
)

// validates a value against a parameter
func (this validate) Param(
	rawValue *model.Data,
	param *model.Param,
) (errs []error) {
	if nil == param {
		return []error{errors.New("Param required")}
	}

	switch {
	case nil != param.Dir:
		var value *string
		if nil != rawValue {
			value = rawValue.Dir
		}
		errs = this.dirParam(value, param.Dir)
	case nil != param.File:
		var value *string
		if nil != rawValue {
			value = rawValue.File
		}
		errs = this.fileParam(value, param.File)
	case nil != param.String:
		var value *string
		if nil != rawValue {
			value = rawValue.String
		}
		errs = this.stringParam(value, param.String)
	case nil != param.Number:
		var value *float64
		if nil != rawValue {
			value = rawValue.Number
		}
		errs = this.numberParam(value, param.Number)
	case nil != param.Socket:
		var value *string
		if nil != rawValue {
			value = rawValue.Socket
		}
		errs = this.socketParam(value, param.Socket)
	}
	return
}
