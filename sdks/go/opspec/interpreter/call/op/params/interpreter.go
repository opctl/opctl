package params

import (
	"github.com/opctl/opctl/sdks/go/model"
)

// Interpret params
func Interpret(
	scope map[string]*model.Value,
	params map[string]*model.ParamSpec,
	opPath string,
	opScratchDir string,
) (
	map[string]*model.Value,
	error,
) {

	var err error
	scope, err = Coerce(scope, params, opScratchDir)
	if err != nil {
		return scope, err
	}

	argsWithDefaults, err := ApplyDefaults(scope, params, opPath, opScratchDir)
	if err != nil {
		return nil, err
	}

	return argsWithDefaults, Validate(argsWithDefaults, params)

}
