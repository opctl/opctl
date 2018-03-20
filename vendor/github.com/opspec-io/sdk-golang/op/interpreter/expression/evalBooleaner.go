package expression

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression/interpolater"
)

type evalBooleaner interface {
	// EvalToBoolean evaluates an expression to an boolean value
	// expression must be a type supported by coerce.ToBoolean
	EvalToBoolean(
		scope map[string]*model.Value,
		expression interface{},
		opDirHandle model.DataHandle,
	) (*model.Value, error)
}

func newEvalBooleaner() evalBooleaner {
	return _evalBooleaner{
		coerce:       coerce.New(),
		interpolater: interpolater.New(),
	}
}

type _evalBooleaner struct {
	coerce       coerce.Coerce
	interpolater interpolater.Interpolater
}

func (ea _evalBooleaner) EvalToBoolean(
	scope map[string]*model.Value,
	expression interface{},
	opDirHandle model.DataHandle,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case bool:
		return &model.Value{Boolean: &expression}, nil
	case string:
		var value *model.Value
		if ref, ok := tryResolveExplicitRef(expression, scope); ok {
			value = ref
		} else {
			stringValue, err := ea.interpolater.Interpolate(
				expression,
				scope,
				opDirHandle,
			)
			if nil != err {
				return nil, err
			}
			value = &model.Value{String: &stringValue}
		}
		return ea.coerce.ToBoolean(value)
	}
	return nil, fmt.Errorf("unable to evaluate %+v to boolean; unsupported type", expression)
}
