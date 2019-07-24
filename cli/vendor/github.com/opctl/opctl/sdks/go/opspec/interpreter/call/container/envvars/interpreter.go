package envvars

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/object"
	stringpkg "github.com/opctl/opctl/sdks/go/opspec/interpreter/string"
	"github.com/opctl/opctl/sdks/go/types"
)

type Interpreter interface {
	Interpret(
		scope map[string]*types.Value,
		scgContainerCallEnvVars interface{},
		opHandle types.DataHandle,
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
	scope map[string]*types.Value,
	scgContainerCallEnvVars interface{},
	opHandle types.DataHandle,
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
