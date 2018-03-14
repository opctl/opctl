package expression

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type evalBooleaner interface {
	// EvalToBoolean evaluates an expression to an boolean value
	// expression must be a type supported by data.CoerceToBoolean
	EvalToBoolean(
		scope map[string]*model.Value,
		expression interface{},
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func newEvalBooleaner() evalBooleaner {
	return _evalBooleaner{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _evalBooleaner struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (ea _evalBooleaner) EvalToBoolean(
	scope map[string]*model.Value,
	expression interface{},
	pkgHandle model.PkgHandle,
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
				pkgHandle,
			)
			if nil != err {
				return nil, err
			}
			value = &model.Value{String: &stringValue}
		}
		return ea.data.CoerceToBoolean(value)
	}
	return nil, fmt.Errorf("unable to evaluate %+v to boolean; unsupported type", expression)
}
