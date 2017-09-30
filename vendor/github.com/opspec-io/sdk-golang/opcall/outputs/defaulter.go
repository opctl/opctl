package outputs

//go:generate counterfeiter -o ./fakeDefaulter.go --fake-name fakeDefaulter ./ defaulter

import (
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"strings"
)

type defaulter interface {
	// Default returns a map consisting of outputArgs & defaults from outputParams
	Default(
		outputArgs map[string]*model.Value,
		outputParams map[string]*model.Param,
		pkgPath string,
	) map[string]*model.Value
}

func newDefaulter() defaulter {
	return _defaulter{}
}

type _defaulter struct {
}

func (dft _defaulter) Default(
	outputArgs map[string]*model.Value,
	outputParams map[string]*model.Param,
	pkgPath string,
) map[string]*model.Value {
	dcgOpCallOutputs := map[string]*model.Value{}

	for paramName, paramValue := range outputParams {
		// apply defaults
		if nil != paramValue {
			switch {
			case nil != paramValue.Array && nil != paramValue.Array.Default:
				dcgOpCallOutputs[paramName] = &model.Value{Array: paramValue.Array.Default}
			case nil != paramValue.Dir && nil != paramValue.Dir.Default && strings.HasPrefix(*paramValue.Dir.Default, "/"):
				dirValue := filepath.Join(pkgPath, *paramValue.Dir.Default)
				dcgOpCallOutputs[paramName] = &model.Value{Dir: &dirValue}
			case nil != paramValue.File && nil != paramValue.File.Default && strings.HasPrefix(*paramValue.File.Default, "/"):
				fileValue := filepath.Join(pkgPath, *paramValue.File.Default)
				dcgOpCallOutputs[paramName] = &model.Value{File: &fileValue}
			case nil != paramValue.Number && nil != paramValue.Number.Default:
				dcgOpCallOutputs[paramName] = &model.Value{Number: paramValue.Number.Default}
			case nil != paramValue.Object && nil != paramValue.Object.Default:
				dcgOpCallOutputs[paramName] = &model.Value{Object: paramValue.Object.Default}
			case nil != paramValue.String && nil != paramValue.String.Default:
				dcgOpCallOutputs[paramName] = &model.Value{String: paramValue.String.Default}
			}
		}
	}
	for argName, argValue := range outputArgs {
		dcgOpCallOutputs[argName] = argValue
	}

	return dcgOpCallOutputs
}
