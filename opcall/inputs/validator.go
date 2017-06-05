package inputs

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ validator

import (
	"errors"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/xeipuuv/gojsonschema"
)

type validator interface {
	Validate(
		value *model.Value,
		param *model.Param,
	) (errors []error)
}

func newValidator() validator {
	// register custom format checkers
	gojsonschema.FormatCheckers.Add("docker-image-ref", DockerImageRefFormatChecker{})
	gojsonschema.FormatCheckers.Add("integer", IntegerFormatChecker{})
	gojsonschema.FormatCheckers.Add("semver", SemVerFormatChecker{})

	return _validator{
		os: ios.New(),
	}
}

type _validator struct {
	os ios.IOS
}

// validates a value against a parameter
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
