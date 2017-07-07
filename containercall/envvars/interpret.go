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
				return nil, fmt.Errorf("Unable to bind env var to '%v' via implicit ref. '%v' is not in scope", envVarName, envVarName)
			}

			switch {
			case nil != value.String:
				dcgContainerCallEnvVars[envVarName] = *value.String
			case nil != value.Number:
				dcgContainerCallEnvVars[envVarName] = strconv.FormatFloat(*value.Number, 'f', -1, 64)
			}
			continue
		}

		// otherwise interpolate value
		dcgContainerCallEnvVars[envVarName] = ev.interpolater.Interpolate(scgContainerEnvVar, scope)
	}
	return dcgContainerCallEnvVars, nil
}
