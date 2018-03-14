package data

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ validator

import (
	"errors"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/data/coerce"
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
		os:     ios.New(),
		coerce: coerce.New(),
	}
}

type _validator struct {
	coerce coerce.Coerce
	os     ios.IOS
}

// Validate validates a value against a parameter
// note: param defaults aren't considered
func (vdt _validator) Validate(
	value *model.Value,
	param *model.Param,
) (errs []error) {
	if nil == param {
		return []error{errors.New("param required")}
	}

	switch {
	case nil != param.Array:
		errs = vdt.validateArray(value, param.Array.Constraints)
	case nil != param.Dir:
		errs = vdt.validateDir(value)
	case nil != param.File:
		errs = vdt.validateFile(value)
	case nil != param.String:
		errs = vdt.validateString(value, param.String.Constraints)
	case nil != param.Number:
		errs = vdt.validateNumber(value, param.Number.Constraints)
	case nil != param.Object:
		errs = vdt.validateObject(value, param.Object.Constraints)
	case nil != param.Socket:
		errs = vdt.validateSocket(value)
	}
	return
}
