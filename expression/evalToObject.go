package expression

import (
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type evalToObject interface {
	// EvalToObject evaluates an expression to a object value
	EvalToObject(
		scope map[string]*model.Value,
		expression string,
		pkgHandle model.PkgHandle,
	) (map[string]interface{}, error)
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
	expression string,
	pkgHandle model.PkgHandle,
) (map[string]interface{}, error) {
	value, err := itp.interpolater.Interpolate(
		expression,
		scope,
		pkgHandle,
	)

	if nil != err {
		return nil, err
	}

	return itp.data.CoerceToObject(value)
}
