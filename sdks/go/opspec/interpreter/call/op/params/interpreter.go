package params

import (
	"github.com/opctl/opctl/sdks/go/model"
)

//Interpret params
func Interpret(
	scope map[string]*model.Value,
	params map[string]*model.Param,
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

	argsWithDefaults := ApplyDefaults(scope, params, opPath)

	return argsWithDefaults, Validate(argsWithDefaults, params)

}
