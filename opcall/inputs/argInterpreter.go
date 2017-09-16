package inputs

//go:generate counterfeiter -o ./fakeArgInterpreter.go --fake-name fakeArgInterpreter ./ argInterpreter

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/expression"
	"github.com/opspec-io/sdk-golang/model"
	"strconv"
)

type argInterpreter interface {
	Interpret(
		name string,
		value interface{},
		param *model.Param,
		parentPkgHandle model.PkgHandle,
		scope map[string]*model.Value,
		opScratchDir string,
	) (*model.Value, error)
}

func newArgInterpreter() argInterpreter {
	return _argInterpreter{
		expression: expression.New(),
	}
}

type _argInterpreter struct {
	expression expression.Expression
}

func (ai _argInterpreter) Interpret(
	name string,
	value interface{},
	param *model.Param,
	parentPkgHandle model.PkgHandle,
	scope map[string]*model.Value,
	opScratchDir string,
) (*model.Value, error) {

	if nil == param {
		return nil, fmt.Errorf("unable to bind to '%v'; '%v' not a defined input", name, name)
	}

	if nil == value || "" == value {
		// implicitly bound
		dcgValue, ok := scope[name]
		if !ok {
			return nil, fmt.Errorf("unable to bind to '%v' via implicit ref; '%v' not in scope", name, name)
		}
		return dcgValue, nil
	}

	// @TODO: move all others to this switch
	switch {
	case nil != param.File:
		fileValue, err := ai.expression.EvalToFile(scope, value, parentPkgHandle, opScratchDir)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%#v'; error was: '%v'", name, value, err.Error())
		}
		return fileValue, nil
	}

	if numberValue, ok := value.(float64); ok {
		switch {
		case nil != param.String:
			stringValue := strconv.FormatFloat(numberValue, 'f', -1, 64)
			return &model.Value{String: &stringValue}, nil
		case nil != param.Number:
			return &model.Value{Number: &numberValue}, nil
		}
	}

	if objectValue, ok := value.(map[string]interface{}); ok {
		return &model.Value{Object: objectValue}, nil
	}

	if stringValue, ok := value.(string); ok {

		if dcgValue, ok := scope[stringValue]; ok {
			// deprecated explicit ref
			return dcgValue, nil
		}

		switch {
		case nil != param.Dir:
			dirValue, err := ai.expression.EvalToDir(scope, stringValue, parentPkgHandle)
			if nil != err {
				return nil, fmt.Errorf("unable to bind '%v' to '%v'; error was: '%v'", name, stringValue, err.Error())
			}
			return dirValue, nil
		case nil != param.Number:
			numberValue, err := ai.expression.EvalToNumber(scope, stringValue, parentPkgHandle)
			if nil != err {
				return nil, fmt.Errorf("unable to bind '%v' to '%v'; error was: '%v'", name, stringValue, err.Error())
			}
			return &model.Value{Number: &numberValue}, nil
		case nil != param.Object:
			objectValue, err := ai.expression.EvalToObject(scope, stringValue, parentPkgHandle)
			if nil != err {
				return nil, fmt.Errorf("unable to bind '%v' to '%v'; error was: '%v'", name, stringValue, err.Error())
			}
			return &model.Value{Object: objectValue}, nil
		case nil != param.String:
			stringValue, err := ai.expression.EvalToString(scope, stringValue, parentPkgHandle)
			if nil != err {
				return nil, fmt.Errorf("unable to bind '%v' to '%v'; error was: '%v'", name, stringValue, err.Error())
			}
			return &model.Value{String: &stringValue}, nil
		case nil != param.Socket:
			return nil, fmt.Errorf("unable to bind '%v' to '%v'; sockets must be passed by reference", name, stringValue)
		}
	}

	return nil, fmt.Errorf("unable to bind '%v' to '%v'", name, value)
}
