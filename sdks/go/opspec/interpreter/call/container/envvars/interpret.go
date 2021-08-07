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
	if containerCallSpecEnvVars == nil {
		return nil, nil
	}

	envVarsMap, err := object.Interpret(
		scope,
		containerCallSpecEnvVars,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to interpret '%v' as envVars: %w", containerCallSpecEnvVars, err)
	}

	envVarsStringMap := map[string]string{}
	for envVarName, envVarValueInterface := range *envVarsMap.Object {
		envVarValue, err := value.Construct(envVarValueInterface)
		if err != nil {
			return nil, fmt.Errorf("unable to construct value for env var '%s': %w", envVarName, err)
		}

		envVarValueString, err := coerce.ToString(envVarValue)
		if err != nil {
			return nil, fmt.Errorf("unable to interpret '%+v' as value of env var '%v': %w", envVarValue, envVarName, err)
		}

		envVarsStringMap[envVarName] = *envVarValueString.String
	}
	return envVarsStringMap, nil
}
