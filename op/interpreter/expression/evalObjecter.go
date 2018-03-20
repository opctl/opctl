package expression

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression/interpolater"
)

type evalObjecter interface {
	// EvalToObject evaluates an expression to a object value
	// expression must be a type supported by coerce.ToObject
	EvalToObject(
		scope map[string]*model.Value,
		expression interface{},
		opHandle model.DataHandle,
	) (*model.Value, error)
}

func newEvalObjecter() evalObjecter {
	return _evalObjecter{
		evalObjectInitializerer: newEvalObjectInitializerer(),
		coerce:                  coerce.New(),
		interpolater:            interpolater.New(),
	}
}

type _evalObjecter struct {
	evalObjectInitializerer
	coerce       coerce.Coerce
	interpolater interpolater.Interpolater
}

func (eo _evalObjecter) EvalToObject(
	scope map[string]*model.Value,
	expression interface{},
	opHandle model.DataHandle,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case map[string]interface{}:
		objectValue, err := eo.evalObjectInitializerer.Eval(
			expression,
			scope,
			opHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %+v to object; error was %v", expression, err)
		}

		return &model.Value{Object: objectValue}, nil
	case string:
		var value *model.Value
		if ref, ok := tryResolveExplicitRef(expression, scope); ok {
			value = ref
		} else {
			stringValue, err := eo.interpolater.Interpolate(
				expression,
				scope,
				opHandle,
			)
			if nil != err {
				return nil, err
			}
			value = &model.Value{String: &stringValue}
		}
		return eo.coerce.ToObject(value)
	}

	return nil, fmt.Errorf("unable to evaluate %+v to object; unsupported type", expression)
}
