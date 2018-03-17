package envvars

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression"
)

type Interpreter interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallEnvVars map[string]interface{},
		opDirHandle model.DataHandle,
	) (map[string]string, error)
}

// NewInterpreter returns a new Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		expression: expression.New(),
	}
}

type _interpreter struct {
	expression expression.Expression
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	scgContainerCallEnvVars map[string]interface{},
	opDirHandle model.DataHandle,
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

		stringValue, err := itp.expression.EvalToString(scope, envVarExpression, opDirHandle)
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
