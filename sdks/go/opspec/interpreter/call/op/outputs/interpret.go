package outputs

import (
	"fmt"
	"sort"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params"
)

// Interpret applies defaults to & validates output args
func Interpret(
	outputArgs map[string]*model.Value,
	outputParams map[string]*model.Param,
	callOutputs map[string]string,
	opPath string,
	opScratchDir string,
) (
	map[string]*model.Value,
	error,
) {
	outputArgs, err := params.Coerce(outputArgs, outputParams, opScratchDir)
	if err != nil {
		return outputArgs, err
	}

	argsWithDefaults := params.ApplyDefaults(outputArgs, outputParams, opPath)

	// ensure the called op supplies each output the callee is expecting
	for callOutputParamName, callOutputParamBoundName := range callOutputs {
		if _, ok := outputParams[callOutputParamName]; !ok {
			// maybe the user has flipped the key and value?
			if _, ok := outputParams[opspec.RefToName(callOutputParamBoundName)]; ok {
				return nil, fmt.Errorf(
					"unknown output '%s', did you mean to use `%s: %s`?",
					callOutputParamName,
					opspec.RefToName(callOutputParamBoundName),
					opspec.NameToRef(callOutputParamName),
				)
			}

			// try to figure out what the user should have provided
			var missingCallOutputParams []string
			for name := range outputParams {
				if _, ok := callOutputs[name]; !ok {
					missingCallOutputParams = append(missingCallOutputParams, name)
				}
			}
			// different environments have different ordering behaviors
			sort.Strings(missingCallOutputParams)

			if len(missingCallOutputParams) == 1 {
				return nil, fmt.Errorf(
					"unknown output '%s', expected '%s'",
					callOutputParamName,
					missingCallOutputParams[0],
				)
			} else if len(missingCallOutputParams) > 0 {
				return nil, fmt.Errorf(
					"unknown output '%s', expected one of [%s]",
					callOutputParamName,
					strings.Join(missingCallOutputParams, ", "),
				)
			} else {
				return nil, fmt.Errorf("unknown output '%s'", callOutputParamName)
			}
		}
	}

	return argsWithDefaults, params.Validate(argsWithDefaults, outputParams)
}
