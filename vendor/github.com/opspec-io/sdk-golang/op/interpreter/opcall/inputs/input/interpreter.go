package input

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression"
)

type Interpreter interface {
	Interpret(
		name string,
		value interface{},
		param *model.Param,
		parentOpDirHandle model.DataHandle,
		scope map[string]*model.Value,
		opScratchDir string,
	) (*model.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		expression: expression.New(),
	}
}

type _interpreter struct {
	expression expression.Expression
}

func (itp _interpreter) Interpret(
	name string,
	valueExpression interface{},
	param *model.Param,
	parentOpDirHandle model.DataHandle,
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
		arrayValue, err := itp.expression.EvalToArray(scope, valueExpression, parentOpDirHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return arrayValue, nil
	case nil != param.Boolean:
		booleanValue, err := itp.expression.EvalToBoolean(scope, valueExpression, parentOpDirHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return booleanValue, nil
	case nil != param.File:
		fileValue, err := itp.expression.EvalToFile(scope, valueExpression, parentOpDirHandle, opScratchDir)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return fileValue, nil
	case nil != param.Dir:
		dirValue, err := itp.expression.EvalToDir(scope, valueExpression, parentOpDirHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return dirValue, nil
	case nil != param.Number:
		numberValue, err := itp.expression.EvalToNumber(scope, valueExpression, parentOpDirHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return numberValue, nil
	case nil != param.Object:
		objectValue, err := itp.expression.EvalToObject(scope, valueExpression, parentOpDirHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return objectValue, nil
	case nil != param.String:
		stringValue, err := itp.expression.EvalToString(scope, valueExpression, parentOpDirHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return stringValue, nil
	case nil != param.Socket:
		return nil, fmt.Errorf("unable to bind '%v' to '%+v'; sockets must be passed by reference", name, valueExpression)
	}

	return nil, fmt.Errorf("unable to bind '%v' to '%v'", name, valueExpression)
}
