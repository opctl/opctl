package expression

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type objectEvaluator interface {
	// EvalToObject evaluates an expression to a object value
	// expression must be a type supported by data.CoerceToObject
	EvalToObject(
		scope map[string]*model.Value,
		expression interface{},
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func newObjectEvaluator() objectEvaluator {
	return _objectEvaluator{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _objectEvaluator struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (eto _objectEvaluator) EvalToObject(
	scope map[string]*model.Value,
	expression interface{},
	pkgHandle model.PkgHandle,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case map[string]interface{}:
		return &model.Value{Object: expression}, nil
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
		return eto.data.CoerceToObject(value)
	}

	return nil, fmt.Errorf("unable to evaluate %+v to object; unsupported type", expression)
}
