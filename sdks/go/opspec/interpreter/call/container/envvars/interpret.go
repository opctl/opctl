package envvars

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/object"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"
)

// Interpret container envVars
func Interpret(
	scope map[string]*model.Value,
	containerCallSpecEnvVars interface{},
) (map[string]string, error) {
	if nil == containerCallSpecEnvVars {
		return nil, nil
	}

	envVarsMap, err := object.Interpret(
		scope,
		containerCallSpecEnvVars,
	)
	if nil != err {
		return nil, fmt.Errorf(
			"unable to interpret '%v' as envVars; error was %v",
			containerCallSpecEnvVars,
			err,
		)
	}

	envVarsStringMap := map[string]string{}
	for envVarName, envVarValueInterface := range *envVarsMap.Object {
		envVarValue, err := value.Construct(envVarValueInterface)
		if nil != err {
			return nil, fmt.Errorf("unable to construct value for env var %s; error was %v", envVarName, err.Error())
		}

		envVarValueString, err := coerce.ToString(envVarValue)
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
