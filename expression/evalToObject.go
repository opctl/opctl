package expression

import (
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
	case string:
		var err error
		if value, err = itp.interpolater.Interpolate(
			expression,
			scope,
			pkgHandle,
		); nil != err {
			return nil, err
		}
	}

	return itp.data.CoerceToObject(value)
}
