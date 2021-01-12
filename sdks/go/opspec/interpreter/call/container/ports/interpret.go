package ports

import (
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
)

// Interpret container ports
func Interpret(
	scope map[string]*model.Value,
	containerCallSpecPorts map[string]string,
	scratchDirPath string,
) (map[string]string, error) {
	containerCallPorts := map[string]string{}
	for callSpecContainerPortNumber, callSpecContainerPortExpression := range containerCallSpecPorts {
		// @TODO: use reference.interpret once reference syntax no longer optional
		callSpecContainerPortExpression = strings.TrimSuffix(strings.TrimPrefix(callSpecContainerPortExpression, reference.RefStart), reference.RefEnd)

		if boundArg, ok := scope[callSpecContainerPortExpression]; ok {
			// bound to var
			containerCallPorts[callSpecContainerPortNumber] = *boundArg.Port
		} else {
			// TODO
		}
	}
	return containerCallPorts, nil
}
