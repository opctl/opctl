package expression

import (
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type evalToDir interface {
	// EvalToDir evaluates an expression to a number value
	EvalToDir(
		scope map[string]*model.Value,
		expression string,
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func newEvalToDir() evalToDir {
	return _evalToDir{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _evalToDir struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (itp _evalToDir) EvalToDir(
	scope map[string]*model.Value,
	expression string,
	pkgHandle model.PkgHandle,
) (*model.Value, error) {
	return itp.interpolater.Interpolate(
		expression,
		scope,
		pkgHandle,
	)
}
