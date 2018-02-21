package expression

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type evalObjecter interface {
	// EvalToObject evaluates an expression to a object value
	// expression must be a type supported by data.CoerceToObject
	EvalToObject(
		scope map[string]*model.Value,
		expression interface{},
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func newEvalObjecter() evalObjecter {
	return _evalObjecter{
		evalObjectInitializerer: newEvalObjectInitializerer(),
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _evalObjecter struct {
	evalObjectInitializerer
	data         data.Data
	interpolater interpolater.Interpolater
}

func (eo _evalObjecter) EvalToObject(
	scope map[string]*model.Value,
	expression interface{},
	pkgHandle model.PkgHandle,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case map[string]interface{}:
		objectValue, err := eo.evalObjectInitializerer.Eval(
			expression,
			scope,
			pkgHandle,
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
				pkgHandle,
			)
			if nil != err {
				return nil, err
			}
			value = &model.Value{String: &stringValue}
		}
		return eo.data.CoerceToObject(value)
	}

	return nil, fmt.Errorf("unable to evaluate %+v to object; unsupported type", expression)
}
