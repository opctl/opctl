package inputs

import (
	"bytes"
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/inputs/input"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/sensitivity"
)

// Interpret interprets inputs via the provided inputArgs, inputParams, and scope;
// opScratchDir will be used to store any run data such as type coercions to files
func Interpret(
	inputArgs map[string]interface{},
	inputParams map[string]*model.ParamSpec,
	opPath string,
	scope map[string]*model.Value,
	opScratchDir string,
	privacy ...map[string]bool,
) (map[string]*model.Value, error) {
	protected := map[string]bool{}
	if len(privacy) > 0 {
		protected = privacy[0]
	}
	interpretedArgs := map[string]*model.Value{}

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
			if protected[argName] {
				err = fmt.Errorf("sensitive input %q failed interpretation", argName)
			} else {
				fmt.Println(argName, err.Error())
			}
			paramErrMap[argName] = err
		}
	}

	for name, value := range interpretedArgs {
		interpretedArgs[name] = sensitivity.Mark(value, protected[name])
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

	// 2) apply protected defaults before the legacy path so their errors stay private.
	for name, param := range inputParams {
		if !protected[name] {
			continue
		}
		var err error
		interpretedArgs, err = params.ApplyDefaults(interpretedArgs, map[string]*model.ParamSpec{name: param}, opPath, opScratchDir)
		if err != nil {
			return nil, fmt.Errorf("unable to interpret input defaults: sensitive parameter %q failed default interpretation", name)
		}
	}
	argsWithDefaults, err := params.ApplyDefaults(interpretedArgs, inputParams, opPath, opScratchDir)
	if err != nil {
		return nil, fmt.Errorf("unable to interpret input defaults: %w", err)
	}

	for name, value := range argsWithDefaults {
		argsWithDefaults[name] = sensitivity.Mark(value, protected[name])
	}
	// 3) validate
	return argsWithDefaults, params.Validate(argsWithDefaults, inputParams)
}
