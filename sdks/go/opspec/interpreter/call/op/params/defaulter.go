package params

//go:generate counterfeiter -o ./fakeDefaulter.go --fake-name FakeDefaulter ./ Defaulter

import (
	"github.com/opctl/opctl/sdks/go/types"
	"path/filepath"
	"strings"
)

// Defaulter defaults params by applying defaults to those w/out corresponding args
type Defaulter interface {
	// Default returns a map consisting of args & defaults from params
	Default(
		args map[string]*types.Value,
		params map[string]*types.Param,
		opPath string,
	) map[string]*types.Value
}

// NewDefaulter returns an initialized Defaulter instance
func NewDefaulter() Defaulter {
	return _defaulter{}
}

type _defaulter struct {
}

func (dft _defaulter) Default(
	args map[string]*types.Value,
	params map[string]*types.Param,
	opPath string,
) map[string]*types.Value {
	argsWithDefaults := map[string]*types.Value{}

	// 1) default all params
	for paramName, paramValue := range params {
		// apply defaults
		if nil != paramValue {
			switch {
			case nil != paramValue.Array && nil != paramValue.Array.Default:
				argsWithDefaults[paramName] = &types.Value{Array: paramValue.Array.Default}
			case nil != paramValue.Boolean && nil != paramValue.Boolean.Default:
				argsWithDefaults[paramName] = &types.Value{Boolean: paramValue.Boolean.Default}
			case nil != paramValue.Dir && nil != paramValue.Dir.Default && strings.HasPrefix(*paramValue.Dir.Default, "/"):
				dirValue := filepath.Join(opPath, *paramValue.Dir.Default)
				argsWithDefaults[paramName] = &types.Value{Dir: &dirValue}
			case nil != paramValue.File && nil != paramValue.File.Default && strings.HasPrefix(*paramValue.File.Default, "/"):
				fileValue := filepath.Join(opPath, *paramValue.File.Default)
				argsWithDefaults[paramName] = &types.Value{File: &fileValue}
			case nil != paramValue.Number && nil != paramValue.Number.Default:
				argsWithDefaults[paramName] = &types.Value{Number: paramValue.Number.Default}
			case nil != paramValue.Object && nil != paramValue.Object.Default:
				argsWithDefaults[paramName] = &types.Value{Object: paramValue.Object.Default}
			case nil != paramValue.String && nil != paramValue.String.Default:
				argsWithDefaults[paramName] = &types.Value{String: paramValue.String.Default}
			}
		}
	}

	// 2) override defaults w/ values (if provided)
	for argName, argValue := range args {
		argsWithDefaults[argName] = argValue
	}

	return argsWithDefaults
}
