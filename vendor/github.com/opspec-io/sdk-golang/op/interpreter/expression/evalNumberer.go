package expression

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression/interpolater"
)

type evalNumberer interface {
	// EvalToNumber evaluates an expression to a number value
	// expression must be a type supported by coerce.ToNumber
	EvalToNumber(
		scope map[string]*model.Value,
		expression interface{},
		opDirHandle model.DataHandle,
	) (*model.Value, error)
}

func newEvalNumberer() evalNumberer {
	return _evalNumberer{
		coerce:       coerce.New(),
		interpolater: interpolater.New(),
	}
}

type _evalNumberer struct {
	coerce       coerce.Coerce
	interpolater interpolater.Interpolater
}

func (en _evalNumberer) EvalToNumber(
	scope map[string]*model.Value,
	expression interface{},
	opDirHandle model.DataHandle,
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
				opDirHandle,
			)
			if nil != err {
				return nil, err
			}

			value = &model.Value{String: &stringValue}
		}
		return en.coerce.ToNumber(value)
	}

	return nil, fmt.Errorf("unable to evaluate %+v to number; unsupported type", expression)
}
