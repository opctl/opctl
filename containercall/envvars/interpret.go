package envvars

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"strconv"
)

func (ev _EnvVars) Interpret(
	scope map[string]*model.Value,
	scgContainerCallEnvVars map[string]string,
) (map[string]string, error) {
	dcgContainerCallEnvVars := map[string]string{}
	for envVarName, scgContainerEnvVar := range scgContainerCallEnvVars {
		if "" == scgContainerEnvVar {
			// implicitly bound
			value, ok := scope[envVarName]
			if !ok {
				return nil, fmt.Errorf(
					"Unable to bind env var to '%v' via implicit ref; '%v' not in scope",
					envVarName,
					envVarName,
				)
			}

			switch {
			case nil != value.Number:
				dcgContainerCallEnvVars[envVarName] = strconv.FormatFloat(*value.Number, 'f', -1, 64)
			case nil != value.Object:
				objectBytes, err := ev.json.Marshal(value.Object)
				if nil != err {
					return nil, fmt.Errorf(
						"Unable to bind env var '%v' via implicit ref; error was: %v",
						envVarName,
						err.Error(),
					)
				}
				dcgContainerCallEnvVars[envVarName] = string(objectBytes)
			case nil != value.String:
				dcgContainerCallEnvVars[envVarName] = *value.String
			}
			continue
		}

		// otherwise interpolate value
		dcgContainerCallEnvVars[envVarName] = ev.interpolater.Interpolate(scgContainerEnvVar, scope)
	}
	return dcgContainerCallEnvVars, nil
}
