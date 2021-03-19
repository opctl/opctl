package params

import (
	"github.com/opctl/opctl/sdks/go/model"
	"path/filepath"
	"strings"
)

// ApplyDefaults to params w/out corresponding args
func ApplyDefaults(
	args map[string]*model.Value,
	params map[string]*model.Param,
	opPath string,
) map[string]*model.Value {
	argsWithDefaults := map[string]*model.Value{}

	// 1) default all params
	for paramName, paramValue := range params {
		// apply defaults
		if paramValue != nil {
			switch {
			case paramValue.Array != nil && paramValue.Array.Default != nil:
				argsWithDefaults[paramName] = &model.Value{Array: paramValue.Array.Default}
			case paramValue.Boolean != nil && paramValue.Boolean.Default != nil:
				argsWithDefaults[paramName] = &model.Value{Boolean: paramValue.Boolean.Default}
			case paramValue.Dir != nil && paramValue.Dir.Default != nil && strings.HasPrefix(*paramValue.Dir.Default, "/"):
				dirValue := filepath.Join(opPath, *paramValue.Dir.Default)
				argsWithDefaults[paramName] = &model.Value{Dir: &dirValue}
			case paramValue.File != nil && paramValue.File.Default != nil && strings.HasPrefix(*paramValue.File.Default, "/"):
				fileValue := filepath.Join(opPath, *paramValue.File.Default)
				argsWithDefaults[paramName] = &model.Value{File: &fileValue}
			case paramValue.Number != nil && paramValue.Number.Default != nil:
				argsWithDefaults[paramName] = &model.Value{Number: paramValue.Number.Default}
			case paramValue.Object != nil && paramValue.Object.Default != nil:
				argsWithDefaults[paramName] = &model.Value{Object: paramValue.Object.Default}
			case paramValue.String != nil && paramValue.String.Default != nil:
				argsWithDefaults[paramName] = &model.Value{String: paramValue.String.Default}
			}
		}
	}

	// 2) override defaults w/ values (if provided)
	for argName, argValue := range args {
		argsWithDefaults[argName] = argValue
	}

	return argsWithDefaults
}
