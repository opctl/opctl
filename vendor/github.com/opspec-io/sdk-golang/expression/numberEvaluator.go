package expression

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type numberEvaluator interface {
	// EvalToNumber evaluates an expression to a number value
	// expression must be a type supported by data.CoerceToNumber
	EvalToNumber(
		scope map[string]*model.Value,
		expression interface{},
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func newNumberEvaluator() numberEvaluator {
	return _numberEvaluator{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _numberEvaluator struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (etn _numberEvaluator) EvalToNumber(
	scope map[string]*model.Value,
	expression interface{},
	pkgHandle model.PkgHandle,
) (*model.Value, error) {
	var value *model.Value

	switch expression := expression.(type) {
	case float64:
		return &model.Value{Number: &expression}, nil
	case map[string]interface{}:
		value = &model.Value{Object: expression}
	case string:
		if ref, ok := tryResolveExplicitRef(expression, scope); ok {
			value = ref
		} else {
			stringValue, err := etn.interpolater.Interpolate(
				expression,
				scope,
				pkgHandle,
			)
			if nil != err {
				return nil, err
			}

			value = &model.Value{String: &stringValue}
		}
	default:
		return nil, fmt.Errorf("unable to evaluate %+v to number; unsupported type", expression)
	}

	return etn.data.CoerceToNumber(value)
}
