package params

//go:generate counterfeiter -o ./fakeDefaulter.go --fake-name FakeDefaulter ./ Defaulter

import (
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"strings"
)

// Defaulter defaults params by applying defaults to those w/out corresponding args
type defaulter interface {
	// Default returns a map consisting of args & defaults from params
	Default(
		args map[string]*model.Value,
		params map[string]*model.Param,
		pkgPath string,
	) map[string]*model.Value
}

// newDefaulter returns a new defaulter
func newDefaulter() defaulter {
	return _defaulter{}
}

type _defaulter struct {
}

func (dft _defaulter) Default(
	args map[string]*model.Value,
	params map[string]*model.Param,
	pkgPath string,
) map[string]*model.Value {
	argsWithDefaults := map[string]*model.Value{}

	for paramName, paramValue := range params {
		// apply defaults
		if nil != paramValue {
			switch {
			case nil != paramValue.Array && nil != paramValue.Array.Default:
				argsWithDefaults[paramName] = &model.Value{Array: paramValue.Array.Default}
			case nil != paramValue.Boolean && nil != paramValue.Boolean.Default:
				argsWithDefaults[paramName] = &model.Value{Boolean: paramValue.Boolean.Default}
			case nil != paramValue.Dir && nil != paramValue.Dir.Default && strings.HasPrefix(*paramValue.Dir.Default, "/"):
				dirValue := filepath.Join(pkgPath, *paramValue.Dir.Default)
				argsWithDefaults[paramName] = &model.Value{Dir: &dirValue}
			case nil != paramValue.File && nil != paramValue.File.Default && strings.HasPrefix(*paramValue.File.Default, "/"):
				fileValue := filepath.Join(pkgPath, *paramValue.File.Default)
				argsWithDefaults[paramName] = &model.Value{File: &fileValue}
			case nil != paramValue.Number && nil != paramValue.Number.Default:
				argsWithDefaults[paramName] = &model.Value{Number: paramValue.Number.Default}
			case nil != paramValue.Object && nil != paramValue.Object.Default:
				argsWithDefaults[paramName] = &model.Value{Object: paramValue.Object.Default}
			case nil != paramValue.String && nil != paramValue.String.Default:
				argsWithDefaults[paramName] = &model.Value{String: paramValue.String.Default}
			}
		}
	}
	for argName, argValue := range args {
		argsWithDefaults[argName] = argValue
	}

	return argsWithDefaults
}
