package envvars

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"

	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/object"
	stringpkg "github.com/opctl/sdk-golang/opspec/interpreter/string"
)

type Interpreter interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallEnvVars interface{},
		opHandle model.DataHandle,
	) (map[string]string, error)
}

// NewInterpreter returns a new Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		objectInterpreter: object.NewInterpreter(),
		stringInterpreter: stringpkg.NewInterpreter(),
		coerce:            coerce.New(),
	}
}

type _interpreter struct {
	coerce            coerce.Coerce
	objectInterpreter object.Interpreter
	stringInterpreter stringpkg.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	scgContainerCallEnvVars interface{},
	opHandle model.DataHandle,
) (map[string]string, error) {
	if nil == scgContainerCallEnvVars {
		return nil, nil
	}

	envVarsMap, err := itp.objectInterpreter.Interpret(
		scope,
		scgContainerCallEnvVars,
		opHandle,
	)
	if nil != err {
		return nil, fmt.Errorf(
			"unable to interpret '%v' as envVars; error was %v",
			scgContainerCallEnvVars,
			err,
		)
	}

	envVarsStringMap := map[string]string{}
	for envVarName, envVarValue := range *envVarsMap.Object {
		envVarValueString, err := itp.stringInterpreter.Interpret(scope, envVarValue, opHandle)
		if nil != err {
			return nil, fmt.Errorf(
				"unable to interpret %+v as value of env var '%v'; error was %v",
				envVarValue,
				envVarName,
				err,
			)
		}

		envVarsStringMap[envVarName] = *envVarValueString.String
	}

	return envVarsStringMap, nil
}
