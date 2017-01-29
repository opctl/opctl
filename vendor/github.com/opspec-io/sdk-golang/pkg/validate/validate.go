package validate

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Validate

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
)

type Validate interface {
	// validates an arg against a parameter
	Param(
		arg *model.Data,
		param *model.Param,
	) (errors []error)
}

func New() Validate {
	return &validate{}
}

type validate struct {
}
