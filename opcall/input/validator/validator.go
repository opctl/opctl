// Package validator validates inputs for a DCG op call
package validator

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Validator

import (
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/xeipuuv/gojsonschema"
)

type Validator interface {
	Validate(
		value *model.Data,
		param *model.Param,
	) (errors []error)
}

func New() Validator {
	// register custom format checkers
	gojsonschema.FormatCheckers.Add("docker-image-ref", DockerImageRefFormatChecker{})
	gojsonschema.FormatCheckers.Add("integer", IntegerFormatChecker{})
	gojsonschema.FormatCheckers.Add("semver", SemVerFormatChecker{})

	return _Validator{
		os: ios.New(),
	}
}

type _Validator struct {
	os ios.IOS
}
