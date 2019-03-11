package initializer

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"github.com/golang-interfaces/gopkg.in-yaml.v2"
	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"
)

type Interpreter interface {
	// Interpret interprets an array initializer expression
	Interpret(
		expression []interface{},
		scope map[string]*model.Value,
		opHandle model.DataHandle,
	) ([]interface{}, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		coerce:       coerce.New(),
		interpolater: interpolater.New(),
		yaml:         iyaml.New(),
	}
}

type _interpreter struct {
	coerce       coerce.Coerce
	interpolater interpolater.Interpolater
	yaml         iyaml.IYAML
}

func (eai _interpreter) Interpret(
	expression []interface{},
	scope map[string]*model.Value,
	opHandle model.DataHandle,
) ([]interface{}, error) {
	arrayBytes, err := eai.yaml.Marshal(expression)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret %+v as array initializer; error was %v", expression, err)
	}

	arrayYAML, err := eai.interpolater.Interpolate(
		string(arrayBytes),
		scope,
		opHandle,
	)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret %+v as array initializer; error was %v", expression, err)
	}

	array := []interface{}{}
	if err := eai.yaml.Unmarshal([]byte(arrayYAML), &array); nil != err {
		return nil, fmt.Errorf("unable to interpret %+v as array initializer; error was %v", expression, err)
	}

	return array, nil
}
