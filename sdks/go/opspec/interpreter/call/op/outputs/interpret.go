package outputs

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params"
)

// Interpret applies defaults to & validates output args
func Interpret(
	outputArgs map[string]*model.Value,
	outputParams map[string]*model.Param,
	opPath string,
	opScratchDir string,
) (
	map[string]*model.Value,
	error,
) {

	outputArgs, err := params.Coerce(outputArgs, outputParams, opScratchDir)
	if nil != err {
		return outputArgs, err
	}

	argsWithDefaults := params.ApplyDefaults(outputArgs, outputParams, opPath)

	return argsWithDefaults, params.Validate(argsWithDefaults, outputParams)

}
