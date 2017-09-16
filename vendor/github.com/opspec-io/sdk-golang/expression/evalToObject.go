package expression

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type evalToObject interface {
	// EvalToObject evaluates an expression to a object value
	// expression must be a type supported by data.CoerceToObject
	EvalToObject(
		scope map[string]*model.Value,
		expression interface{},
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func newEvalToObject() evalToObject {
	return _evalToObject{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _evalToObject struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (itp _evalToObject) EvalToObject(
	scope map[string]*model.Value,
	expression interface{},
	pkgHandle model.PkgHandle,
) (*model.Value, error) {
	var value *model.Value

	switch expression := expression.(type) {
	case float64:
		value = &model.Value{Number: &expression}
	case map[string]interface{}:
		return &model.Value{Object: expression}, nil
	case string:
		var err error
		if value, err = itp.interpolater.Interpolate(
			expression,
			scope,
			pkgHandle,
		); nil != err {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unable to evaluate %+v to object; unsupported type", expression)
	}

	return itp.data.CoerceToObject(value)
}
