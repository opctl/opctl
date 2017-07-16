package envvars

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ EnvVars

import (
	"github.com/golang-interfaces/ijson"
	"github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type EnvVars interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallEnvVars map[string]string,
	) (map[string]string, error)
}

func New() EnvVars {
	return _EnvVars{
		interpolater: interpolater.New(),
		json:         ijson.New(),
	}
}

type _EnvVars struct {
	interpolater interpolater.Interpolater
	json         ijson.IJSON
}
