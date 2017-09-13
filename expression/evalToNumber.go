package expression

import (
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/interpolater"
)

type evalToNumber interface {
	// EvalToNumber evaluates an expression to a number value
	EvalToNumber(
		scope map[string]*model.Value,
		expression string,
		pkgHandle model.PkgHandle,
	) (float64, error)
}

func newEvalToNumber() evalToNumber {
	return _evalToNumber{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _evalToNumber struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (itp _evalToNumber) EvalToNumber(
	scope map[string]*model.Value,
	expression string,
	pkgHandle model.PkgHandle,
) (float64, error) {
	value, err := itp.interpolater.Interpolate(
		expression,
		scope,
		pkgHandle,
	)

	if nil != err {
		return 0, err
	}

	return itp.data.CoerceToNumber(value)
}
