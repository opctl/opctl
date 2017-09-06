package inputs

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ validator

import (
	"errors"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/dir"
	"github.com/opspec-io/sdk-golang/file"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/number"
	"github.com/opspec-io/sdk-golang/object"
	stringPkg "github.com/opspec-io/sdk-golang/string"
)

type validator interface {
	Validate(
		value *model.Value,
		param *model.Param,
	) (errors []error)
}

func newValidator() validator {

	return _validator{
		dir:    dir.New(),
		file:   file.New(),
		os:     ios.New(),
		number: number.New(),
		object: object.New(),
		string: stringPkg.New(),
	}
}

type _validator struct {
	dir    dir.Dir
	file   file.File
	os     ios.IOS
	number number.Number
	object object.Object
	string stringPkg.String
}

// Validate validates a value against a parameter
// note: param defaults aren't considered
func (this _validator) Validate(
	rawValue *model.Value,
	param *model.Param,
) (errs []error) {
	if nil == param {
		return []error{errors.New("param required")}
	}

	switch {
	case nil != param.Dir:
		var value *string
		if nil != rawValue {
			value = rawValue.Dir
		}
		errs = this.dir.Validate(value)
	case nil != param.File:
		var value *string
		if nil != rawValue {
			value = rawValue.File
		}
		errs = this.file.Validate(value)
	case nil != param.String:
		var value *string
		if nil != rawValue {
			value = rawValue.String
		}
		errs = this.string.Validate(value, param.String.Constraints)
	case nil != param.Number:
		var value *float64
		if nil != rawValue {
			value = rawValue.Number
		}
		errs = this.number.Validate(value, param.Number.Constraints)
	case nil != param.Object:
		var value map[string]interface{}
		if nil != rawValue {
			value = rawValue.Object
		}
		errs = this.object.Validate(value, param.Object.Constraints)
	case nil != param.Socket:
		var value *string
		if nil != rawValue {
			value = rawValue.Socket
		}
		errs = this.validateSocket(value)
	}
	return
}
