package params

import (
	"github.com/opctl/opctl/sdks/go/model"
)

// Interpret params
func Interpret(
	scope map[string]*ipld.Node,
	params map[string]*model.ParamSpec,
	opPath string,
	opScratchDir string,
) (
	map[string]*ipld.Node,
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
