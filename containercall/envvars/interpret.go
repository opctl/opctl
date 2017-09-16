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
	for envVarName, envVarExpression := range scgContainerCallEnvVars {
		if "" == envVarExpression {
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

		stringValue, err := ev.expression.EvalToString(scope, envVarExpression, pkgHandle)
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
