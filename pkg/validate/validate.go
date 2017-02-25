package validate

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Validate

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/xeipuuv/gojsonschema"
)

type Validate interface {
	Param(
		value *model.Data,
		param *model.Param,
	) (errors []error)
}

func New() Validate {
	// register custom format checkers
	gojsonschema.FormatCheckers.Add("docker-image-ref", DockerImageRefFormatChecker{})
	gojsonschema.FormatCheckers.Add("integer", IntegerFormatChecker{})
	gojsonschema.FormatCheckers.Add("semver", SemVerFormatChecker{})

	return validate{}
}

type validate struct{}
