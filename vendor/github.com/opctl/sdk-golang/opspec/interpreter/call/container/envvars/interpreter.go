package envvars

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"github.com/opctl/sdk-golang/model"
	stringPkg "github.com/opctl/sdk-golang/opspec/interpreter/string"
)

type Interpreter interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallEnvVars map[string]interface{},
		opHandle model.DataHandle,
	) (map[string]string, error)
}

// NewInterpreter returns a new Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		stringInterpreter: stringPkg.NewInterpreter(),
	}
}

type _interpreter struct {
	stringInterpreter stringPkg.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	scgContainerCallEnvVars map[string]interface{},
	opHandle model.DataHandle,
) (map[string]string, error) {
	dcgContainerCallEnvVars := map[string]string{}
	for envVarName, envVarExpression := range scgContainerCallEnvVars {
		if nil == envVarExpression {
			// implicitly bound
			if _, ok := scope[envVarName]; !ok {
				return nil, fmt.Errorf(
					"unable to bind env var to '%v' via implicit ref; '%v' not in scope",
					envVarName,
					envVarName,
				)
			}
			envVarExpression = fmt.Sprintf("$(%v)", envVarName)
		}

		stringValue, err := itp.stringInterpreter.Interpret(scope, envVarExpression, opHandle)
		if nil != err {
			return nil, fmt.Errorf(
				"unable to bind env var to '%v' via implicit ref; '%v' not in scope",
				envVarName,
				envVarName,
			)
		}

		dcgContainerCallEnvVars[envVarName] = *stringValue.String
	}
	return dcgContainerCallEnvVars, nil
}
