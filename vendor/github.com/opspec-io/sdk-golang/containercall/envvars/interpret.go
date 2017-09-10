package envvars

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
)

func (ev _EnvVars) Interpret(
	scope map[string]*model.Value,
	scgContainerCallEnvVars map[string]string,
	pkgHandle model.PkgHandle,
) (map[string]string, error) {
	dcgContainerCallEnvVars := map[string]string{}
	for envVarName, scgContainerEnvVar := range scgContainerCallEnvVars {
		if "" == scgContainerEnvVar {
			// implicitly bound
			if _, ok := scope[envVarName]; !ok {
				return nil, fmt.Errorf(
					"Unable to bind env var to '%v' via implicit ref; '%v' not in scope",
					envVarName,
					envVarName,
				)
			}
			scgContainerEnvVar = fmt.Sprintf("$(%v)", envVarName)
		}

		stringValue, err := ev.string.Interpret(scope, scgContainerEnvVar, pkgHandle)
		if nil != err {
			return nil, fmt.Errorf(
				"Unable to bind env var to '%v' via implicit ref; '%v' not in scope",
				envVarName,
				envVarName,
			)
		}

		dcgContainerCallEnvVars[envVarName] = stringValue
	}
	return dcgContainerCallEnvVars, nil
}
