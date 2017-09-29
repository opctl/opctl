package envvars

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ EnvVars

import (
	"github.com/opspec-io/sdk-golang/expression"
	"github.com/opspec-io/sdk-golang/model"
)

type EnvVars interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallEnvVars map[string]interface{},
		pkgHandle model.PkgHandle,
	) (map[string]string, error)
}

func New() EnvVars {
	return _EnvVars{
		expression: expression.New(),
	}
}

type _EnvVars struct {
	expression expression.Expression
}
