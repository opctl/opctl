package validator

import (
	"errors"
	"github.com/opspec-io/sdk-golang/model"
)

// validates a value against a parameter
func (this _Validator) Validate(
	rawValue *model.Data,
	param *model.Param,
) (errs []error) {
	if nil == param {
		return []error{errors.New("Validate required")}
	}

	switch {
	case nil != param.Dir:
		var value *string
		if nil != rawValue {
			value = rawValue.Dir
		}
		errs = this.validateDir(value, param.Dir)
	case nil != param.File:
		var value *string
		if nil != rawValue {
			value = rawValue.File
		}
		errs = this.validateFile(value, param.File)
	case nil != param.String:
		var value *string
		if nil != rawValue {
			value = rawValue.String
		}
		errs = this.validateString(value, param.String)
	case nil != param.Number:
		var value *float64
		if nil != rawValue {
			value = rawValue.Number
		}
		errs = this.validateNumber(value, param.Number)
	case nil != param.Socket:
		var value *string
		if nil != rawValue {
			value = rawValue.Socket
		}
		errs = this.validateSocket(value, param.Socket)
	}
	return
}
