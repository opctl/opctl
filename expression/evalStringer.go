package expression

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type evalStringer interface {
	// EvalToString evaluates an expression to a string value
	EvalToString(
		scope map[string]*model.Value,
		expression interface{},
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func newEvalStringer() evalStringer {
	return _evalStringer{
		evalArrayInitializerer:  newEvalArrayInitializerer(),
		evalObjectInitializerer: newEvalObjectInitializerer(),
		coerce:                  coerce.New(),
		interpolater:            interpolater.New(),
	}
}

type _evalStringer struct {
	evalArrayInitializerer
	evalObjectInitializerer
	coerce       coerce.Coerce
	interpolater interpolater.Interpolater
}

func (es _evalStringer) EvalToString(
	scope map[string]*model.Value,
	expression interface{},
	pkgHandle model.PkgHandle,
) (*model.Value, error) {
	var value *model.Value

	switch expression := expression.(type) {
	case float64:
		value = &model.Value{Number: &expression}
	case map[string]interface{}:
		objectValue, err := es.evalObjectInitializerer.Eval(
			expression,
			scope,
			pkgHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %+v to string; error was %v", expression, err)
		}

		value = &model.Value{Object: objectValue}
	case []interface{}:
		arrayValue, err := es.evalArrayInitializerer.Eval(
			expression,
			scope,
			pkgHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %+v to string; error was %v", expression, err)
		}

		value = &model.Value{Array: arrayValue}
	case string:
		if ref, ok := tryResolveExplicitRef(expression, scope); ok {
			value = ref
		} else {
			stringValue, err := es.interpolater.Interpolate(
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
		return nil, fmt.Errorf("unable to evaluate %+v to string; unsupported type", expression)
	}

	return es.coerce.ToString(value)
}
