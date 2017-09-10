package string

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/interpolater"
)

type interpreter interface {
	// interprets an expression to a string
	Interpret(
		scope map[string]*model.Value,
		expression string,
		pkgHandle model.PkgHandle,
	) (string, error)
}

func newInterpreter() interpreter {
	return _interpreter{
		interpolater: interpolater.New(),
	}
}

type _interpreter struct {
	interpolater interpolater.Interpolater
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	expression string,
	pkgHandle model.PkgHandle,
) (string, error) {
	return itp.interpolater.Interpolate(
		expression,
		scope,
		pkgHandle,
	)
}
