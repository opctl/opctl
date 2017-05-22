// Package paramvalidator validates params of a DCG op call
package paramvalidator

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ ParamValidator

import (
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/xeipuuv/gojsonschema"
)

type ParamValidator interface {
	Validate(
		value *model.Data,
		param *model.Param,
	) (errors []error)
}

func New() ParamValidator {
	// register custom format checkers
	gojsonschema.FormatCheckers.Add("docker-image-ref", DockerImageRefFormatChecker{})
	gojsonschema.FormatCheckers.Add("integer", IntegerFormatChecker{})
	gojsonschema.FormatCheckers.Add("semver", SemVerFormatChecker{})

	return _ParamValidator{
		os: ios.New(),
	}
}

type _ParamValidator struct {
	os ios.IOS
}
