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
	// Interpret interprets an object initializer expression
	Interpret(
		expression map[string]interface{},
		scope map[string]*model.Value,
		opHandle model.DataHandle,
	) (map[string]interface{}, error)
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

func (eoi _interpreter) Interpret(
	expression map[string]interface{},
	scope map[string]*model.Value,
	opHandle model.DataHandle,
) (map[string]interface{}, error) {

	// expand shorthand properties w/out mutating original (maps passed by reference in go)
	expressionWithExpandedShorthandProps := map[string]interface{}{}
	for propName, propValue := range expression {
		if nil == propValue {
			expressionWithExpandedShorthandProps[propName] = fmt.Sprintf("%v%v%v", interpolater.RefStart, propName, interpolater.RefEnd)
		} else {
			expressionWithExpandedShorthandProps[propName] = propValue
		}
	}

	objectBytes, err := eoi.yaml.Marshal(expressionWithExpandedShorthandProps)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret %+v as objectInitializer; error was %v", expression, err)
	}

	objectYAML, err := eoi.interpolater.Interpolate(
		string(objectBytes),
		scope,
		opHandle,
	)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret %+v as objectInitializer; error was %v", expression, err)
	}

	object := map[string]interface{}{}
	if err := eoi.yaml.Unmarshal([]byte(objectYAML), &object); nil != err {
		return nil, fmt.Errorf("unable to interpret %+v as objectInitializer; error was %v", expression, err)
	}

	return object, nil
}
