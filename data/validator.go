package data

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ validator

import (
	"errors"
	"github.com/chrisdostert/gojsonschema"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
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
		os:      ios.New(),
		coercer: newCoercer(),
	}
}

type _validator struct {
	coercer coercer
	os      ios.IOS
}

// Validate validates a value against a parameter
// note: param defaults aren't considered
func (this _validator) Validate(
	value *model.Value,
	param *model.Param,
) (errs []error) {
	if nil == param {
		return []error{errors.New("param required")}
	}

	switch {
	case nil != param.Dir:
		errs = this.validateDir(value)
	case nil != param.File:
		errs = this.validateFile(value)
	case nil != param.String:
		errs = this.validateString(value, param.String.Constraints)
	case nil != param.Number:
		errs = this.validateNumber(value, param.Number.Constraints)
	case nil != param.Object:
		errs = this.validateObject(value, param.Object.Constraints)
	case nil != param.Socket:
		errs = this.validateSocket(value)
	}
	return
}
