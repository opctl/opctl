package expression

import (
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type evalToString interface {
	// EvalToString evaluates an expression to a string value
	EvalToString(
		scope map[string]*model.Value,
		expression interface{},
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func newEvalToString() evalToString {
	return _evalToString{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _evalToString struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (itp _evalToString) EvalToString(
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

	return itp.data.CoerceToString(value)
}
