package validate

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Validate

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/virtual-go/fs"
	"github.com/virtual-go/fs/osfs"
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

	return validate{
		fs: osfs.New(),
	}
}

type validate struct {
	fs fs.FS
}
