package expression

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/interpolater"
)

type evalToString interface {
	// Eval evaluates an expression to a string value
	EvalToString(
		scope map[string]*model.Value,
		expression string,
		pkgHandle model.PkgHandle,
	) (string, error)
}

func newEvalToString() evalToString {
	return _evalToString{
		interpolater: interpolater.New(),
	}
}

type _evalToString struct {
	interpolater interpolater.Interpolater
}

func (itp _evalToString) EvalToString(
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
