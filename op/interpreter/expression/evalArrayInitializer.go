package expression

import (
	"fmt"
	"github.com/golang-interfaces/gopkg.in-yaml.v2"
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression/interpolater"
)

//go:generate counterfeiter -o ./fakeEvalArrayInitializerer.go --fake-name fakeEvalArrayInitializerer ./ evalArrayInitializerer

type evalArrayInitializerer interface {
	// Eval evaluates an arrayInitializer expression
	Eval(
		expression []interface{},
		scope map[string]*model.Value,
		opHandle model.DataHandle,
	) ([]interface{}, error)
}

// newEvalArrayInitializerer returns a new evalArrayInitializerer instance
func newEvalArrayInitializerer() evalArrayInitializerer {
	return _evalArrayInitializerer{
		coerce:       coerce.New(),
		interpolater: interpolater.New(),
		yaml:         iyaml.New(),
	}
}

type _evalArrayInitializerer struct {
	coerce       coerce.Coerce
	interpolater interpolater.Interpolater
	yaml         iyaml.IYAML
}

func (eai _evalArrayInitializerer) Eval(
	expression []interface{},
	scope map[string]*model.Value,
	opHandle model.DataHandle,
) ([]interface{}, error) {
	arrayBytes, err := eai.yaml.Marshal(expression)
	if nil != err {
		return nil, fmt.Errorf("unable to eval %+v as arrayInitializer; error was %v", expression, err)
	}

	arrayYAML, err := eai.interpolater.Interpolate(
		string(arrayBytes),
		scope,
		opHandle,
	)
	if nil != err {
		return nil, fmt.Errorf("unable to eval %+v as arrayInitializer; error was %v", expression, err)
	}

	array := []interface{}{}
	if err := eai.yaml.Unmarshal([]byte(arrayYAML), &array); nil != err {
		return nil, fmt.Errorf("unable to eval %+v as arrayInitializer; error was %v", expression, err)
	}

	return array, nil
}
