package expression

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type evalNumberer interface {
	// EvalToNumber evaluates an expression to a number value
	// expression must be a type supported by data.CoerceToNumber
	EvalToNumber(
		scope map[string]*model.Value,
		expression interface{},
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func newEvalNumberer() evalNumberer {
	return _evalNumberer{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _evalNumberer struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (en _evalNumberer) EvalToNumber(
	scope map[string]*model.Value,
	expression interface{},
	pkgHandle model.PkgHandle,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case float64:
		return &model.Value{Number: &expression}, nil
	case string:
		var value *model.Value
		if ref, ok := tryResolveExplicitRef(expression, scope); ok {
			value = ref
		} else {
			stringValue, err := en.interpolater.Interpolate(
				expression,
				scope,
				pkgHandle,
			)
			if nil != err {
				return nil, err
			}

			value = &model.Value{String: &stringValue}
		}
		return en.data.CoerceToNumber(value)
	}

	return nil, fmt.Errorf("unable to evaluate %+v to number; unsupported type", expression)
}
