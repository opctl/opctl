package inputs

//go:generate counterfeiter -o ./fakeArgInterpreter.go --fake-name fakeArgInterpreter ./ argInterpreter

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/expression"
	"github.com/opspec-io/sdk-golang/model"
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
	valueExpression interface{},
	param *model.Param,
	parentPkgHandle model.PkgHandle,
	scope map[string]*model.Value,
	opScratchDir string,
) (*model.Value, error) {

	if nil == param {
		return nil, fmt.Errorf("unable to bind to '%v'; '%v' not a defined input", name, name)
	}

	if nil == valueExpression || "" == valueExpression {
		// implicitly bound
		dcgValue, ok := scope[name]
		if !ok {
			return nil, fmt.Errorf("unable to bind to '%v' via implicit ref; '%v' not in scope", name, name)
		}
		return dcgValue, nil
	}

	switch {
	case nil != param.Array:
		arrayValue, err := ai.expression.EvalToArray(scope, valueExpression, parentPkgHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return arrayValue, nil
	case nil != param.Boolean:
		booleanValue, err := ai.expression.EvalToBoolean(scope, valueExpression, parentPkgHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return booleanValue, nil
	case nil != param.File:
		fileValue, err := ai.expression.EvalToFile(scope, valueExpression, parentPkgHandle, opScratchDir)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return fileValue, nil
	case nil != param.Dir:
		dirValue, err := ai.expression.EvalToDir(scope, valueExpression, parentPkgHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return dirValue, nil
	case nil != param.Number:
		numberValue, err := ai.expression.EvalToNumber(scope, valueExpression, parentPkgHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return numberValue, nil
	case nil != param.Object:
		objectValue, err := ai.expression.EvalToObject(scope, valueExpression, parentPkgHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return objectValue, nil
	case nil != param.String:
		stringValue, err := ai.expression.EvalToString(scope, valueExpression, parentPkgHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return stringValue, nil
	case nil != param.Socket:
		return nil, fmt.Errorf("unable to bind '%v' to '%+v'; sockets must be passed by reference", name, valueExpression)
	}

	return nil, fmt.Errorf("unable to bind '%v' to '%v'", name, valueExpression)
}
