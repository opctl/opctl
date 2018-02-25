package expression

import (
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

//go:generate counterfeiter -o ./fakeEvalArrayInitializerer.go --fake-name fakeEvalArrayInitializerer ./ evalArrayInitializerer

type evalArrayInitializerer interface {
	// Eval evaluates an arrayInitializer expression
	Eval(
		expression []interface{},
		scope map[string]*model.Value,
		pkgHandle model.PkgHandle,
	) ([]interface{}, error)
}

// newEvalArrayInitializerer returns a new evalArrayInitializerer instance
func newEvalArrayInitializerer() evalArrayInitializerer {
	return _evalArrayInitializerer{
		data:         data.New(),
		interpolater: interpolater.New(),
		json:         ijson.New(),
	}
}

type _evalArrayInitializerer struct {
	data         data.Data
	interpolater interpolater.Interpolater
	json         ijson.IJSON
}

func (eoi _evalArrayInitializerer) Eval(
	expression []interface{},
	scope map[string]*model.Value,
	pkgHandle model.PkgHandle,
) ([]interface{}, error) {
	arrayBytes, err := eoi.json.Marshal(expression)
	if nil != err {
		return nil, fmt.Errorf("unable to eval %+v as arrayInitializer; error was %v", expression, err)
	}

	arrayJson, err := eoi.interpolater.Interpolate(
		string(arrayBytes),
		scope,
		pkgHandle,
	)
	if nil != err {
		return nil, fmt.Errorf("unable to eval %+v as arrayInitializer; error was %v", expression, err)
	}

	array := []interface{}{}
	if err := eoi.json.Unmarshal([]byte(arrayJson), &array); nil != err {
		return nil, fmt.Errorf("unable to eval %+v as arrayInitializer; error was %v", expression, err)
	}

	return array, nil
}
