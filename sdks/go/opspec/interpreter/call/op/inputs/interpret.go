package inputs

import (
	"bytes"
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/inputs/input"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params"
)

// Interpret interprets inputs via the provided inputArgs, inputParams, and scope;
// opScratchDir will be used to store any run data such as type coercions to files
func Interpret(
	inputArgs map[string]interface{},
	inputParams map[string]*model.ParamSpec,
	opPath string,
	scope map[string]*ipld.Node,
	opScratchDir string,
) (map[string]*ipld.Node, error) {
	interpretedArgs := map[string]*ipld.Node{}

	// 1) interpret
	paramErrMap := map[string]error{}
	for argName, argValue := range inputArgs {
		var err error
		interpretedArgs[argName], err = input.Interpret(
			argName,
			argValue,
			inputParams[argName],
			scope,
			opScratchDir,
		)
		if err != nil {
			fmt.Println(argName, err.Error())
			paramErrMap[argName] = err
		}
	}

	if len(paramErrMap) > 0 {
		// return error w/ fancy formatted msg
		messageBuffer := bytes.NewBufferString("")
		for paramName, err := range paramErrMap {
			messageBuffer.WriteString(fmt.Sprintf(`
		- %v: %v`, paramName, err.Error()))
		}
		messageBuffer.WriteString(`
`)
		return nil, fmt.Errorf(`
-
  validation error(s):
%v
-`, messageBuffer.String())
	}

	// 2) apply defaults
	argsWithDefaults, err := params.ApplyDefaults(interpretedArgs, inputParams, opPath, opScratchDir)
	if err != nil {
		return nil, fmt.Errorf("unable to interpret input defaults: %w", err)
	}

	// 3) validate
	return argsWithDefaults, params.Validate(argsWithDefaults, inputParams)
}
