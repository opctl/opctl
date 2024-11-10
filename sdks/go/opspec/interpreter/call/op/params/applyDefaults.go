package params

import (
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/array"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/boolean"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/dir"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/file"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/number"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/object"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

// ApplyDefaults to params w/out corresponding args
// opScratchDir will be used to store any run data such as type coercions to files
func ApplyDefaults(
	args map[string]*ipld.Node,
	params map[string]*model.ParamSpec,
	opPath,
	opScratchDir string,
) (map[string]*ipld.Node, error) {

	argsWithDefaults := map[string]*ipld.Node{}
	// don't mutate provided args map
	for argName, argValue := range args {
		argsWithDefaults[argName] = argValue
	}

	parentDirPath := filepath.Dir(opPath)
	defaultsScope := map[string]*ipld.Node{
		// add deprecated absolute path to scope
		"/": {
			Dir: &opPath,
		},
		// add current directory to scope
		"./": {
			Dir: &opPath,
		},
		// add parent directory to scope
		"../": {
			Dir: &parentDirPath,
		},
	}

	for paramName, paramValue := range params {
		// default params which have no corresponding arg
		if _, ok := argsWithDefaults[paramName]; !ok && paramValue != nil {
			switch {
			case paramValue.Array != nil && paramValue.Array.Default != nil:
				defaultValue, err := array.Interpret(
					defaultsScope,
					paramValue.Array.Default,
				)
				if err != nil {
					return nil, err
				}

				argsWithDefaults[paramName] = defaultValue
			case paramValue.Boolean != nil && paramValue.Boolean.Default != nil:
				defaultValue, err := boolean.Interpret(
					defaultsScope,
					paramValue.Boolean.Default,
				)
				if err != nil {
					return nil, err
				}

				argsWithDefaults[paramName] = defaultValue
			case paramValue.Dir != nil && paramValue.Dir.Default != nil:
				defaultExpression := paramValue.Dir.Default

				if defaultExpressionAsString, ok := defaultExpression.(string); ok && strings.HasPrefix(defaultExpressionAsString, "/") {
					// convert deprecated syntax to current syntax
					defaultExpression = opspec.NameToRef(defaultExpressionAsString)
				}

				defaultValue, err := dir.Interpret(
					defaultsScope,
					defaultExpression,
					opScratchDir,
					false,
				)
				if err != nil {
					return nil, err
				}

				argsWithDefaults[paramName] = defaultValue
			case paramValue.File != nil && paramValue.File.Default != nil:
				defaultExpression := paramValue.File.Default

				if defaultExpressionAsString, ok := defaultExpression.(string); ok && strings.HasPrefix(defaultExpressionAsString, "/") {
					// convert deprecated syntax to current syntax
					defaultExpression = opspec.NameToRef(defaultExpressionAsString)
				}

				defaultValue, err := file.Interpret(
					defaultsScope,
					defaultExpression,
					opScratchDir,
					false,
				)
				if err != nil {
					return nil, err
				}

				argsWithDefaults[paramName] = defaultValue
			case paramValue.Number != nil && paramValue.Number.Default != nil:
				defaultValue, err := number.Interpret(
					defaultsScope,
					paramValue.Number.Default,
				)
				if err != nil {
					return nil, err
				}

				argsWithDefaults[paramName] = defaultValue
			case paramValue.Object != nil && paramValue.Object.Default != nil:
				defaultValue, err := object.Interpret(
					defaultsScope,
					paramValue.Object.Default,
				)
				if err != nil {
					return nil, err
				}

				argsWithDefaults[paramName] = defaultValue
			case paramValue.String != nil && paramValue.String.Default != nil:
				defaultValue, err := str.Interpret(
					defaultsScope,
					paramValue.String.Default,
				)
				if err != nil {
					return nil, err
				}

				argsWithDefaults[paramName] = defaultValue
			}
		}
	}

	return argsWithDefaults, nil
}
