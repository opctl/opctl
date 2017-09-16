package expression

import (
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type evalToNumber interface {
	// EvalToNumber evaluates an expression to a number value
	// expression must be a type supported by data.CoerceToNumber
	EvalToNumber(
		scope map[string]*model.Value,
		expression interface{},
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
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
	expression interface{},
	pkgHandle model.PkgHandle,
) (*model.Value, error) {
	var value *model.Value

	switch expression := expression.(type) {
	case string:
		var err error
		if value, err = itp.interpolater.Interpolate(
			expression,
			scope,
			pkgHandle,
		); nil != err {
			return nil, err
		}
	}

	return itp.data.CoerceToNumber(value)
}
