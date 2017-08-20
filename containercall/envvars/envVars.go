package envvars

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ EnvVars

import (
	"github.com/opspec-io/sdk-golang/model"
	stringPkg "github.com/opspec-io/sdk-golang/string"
)

type EnvVars interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallEnvVars map[string]string,
	) (map[string]string, error)
}

func New() EnvVars {
	return _EnvVars{
		string: stringPkg.New(),
	}
}

type _EnvVars struct {
	string stringPkg.String
}
