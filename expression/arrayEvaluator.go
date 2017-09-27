package expression

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type arrayEvaluator interface {
	// EvalToArray evaluates an expression to a array value
	// expression must be a type supported by data.CoerceToArray
	EvalToArray(
		scope map[string]*model.Value,
		expression interface{},
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func newArrayEvaluator() arrayEvaluator {
	return _arrayEvaluator{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _arrayEvaluator struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (eto _arrayEvaluator) EvalToArray(
	scope map[string]*model.Value,
	expression interface{},
	pkgHandle model.PkgHandle,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case []interface{}:
		return &model.Value{Array: expression}, nil
	case string:
		var value *model.Value
		if ref, ok := tryResolveExplicitRef(expression, scope); ok {
			value = ref
		} else {
			stringValue, err := eto.interpolater.Interpolate(
				expression,
				scope,
				pkgHandle,
			)
			if nil != err {
				return nil, err
			}
			value = &model.Value{String: &stringValue}
		}
    return eto.data.CoerceToArray(value)
	}
	return nil, fmt.Errorf("unable to evaluate %+v to array; unsupported type", expression)
}
