package expression

import (
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/interpolater"
)

type evalToString interface {
	// EvalToString evaluates an expression to a string value
	EvalToString(
		scope map[string]*model.Value,
		expression string,
		pkgHandle model.PkgHandle,
	) (string, error)
}

func newEvalToString() evalToString {
	return _evalToString{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _evalToString struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (itp _evalToString) EvalToString(
	scope map[string]*model.Value,
	expression string,
	pkgHandle model.PkgHandle,
) (string, error) {
	value, err := itp.interpolater.Interpolate(
		expression,
		scope,
		pkgHandle,
	)

	if nil != err {
		return "", err
	}

	return itp.data.CoerceToString(value)
}
