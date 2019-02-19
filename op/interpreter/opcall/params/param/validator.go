package param

//go:generate counterfeiter -o ./fakeValidator.go --fake-name FakeValidator ./ Validator

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/xeipuuv/gojsonschema"
)

type Validator interface {
	Validate(
		value *model.Value,
		param *model.Param,
	) []error
}

func NewValidator() Validator {
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
) []error {
	if nil == param {
		return []error{errors.New("param required")}
	}

	switch {
	case nil != param.Array:
		return vdt.validateArray(value, param.Array.Constraints)
	case nil != param.Boolean:
		return vdt.validateBoolean(value)
	case nil != param.Dir:
		return vdt.validateDir(value)
	case nil != param.File:
		return vdt.validateFile(value)
	case nil != param.String:
		return vdt.validateString(value, param.String.Constraints)
	case nil != param.Number:
		return vdt.validateNumber(value, param.Number.Constraints)
	case nil != param.Object:
		return vdt.validateObject(value, param.Object.Constraints)
	case nil != param.Socket:
		return vdt.validateSocket(value)
	default:
		return []error{fmt.Errorf("unable to validate value; param was unexpected type %+v", param)}
	}
}
