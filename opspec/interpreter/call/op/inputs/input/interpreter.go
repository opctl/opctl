package input

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/array"
	"github.com/opctl/sdk-golang/opspec/interpreter/boolean"
	"github.com/opctl/sdk-golang/opspec/interpreter/dir"
	"github.com/opctl/sdk-golang/opspec/interpreter/file"
	"github.com/opctl/sdk-golang/opspec/interpreter/number"
	"github.com/opctl/sdk-golang/opspec/interpreter/object"
	stringPkg "github.com/opctl/sdk-golang/opspec/interpreter/string"
)

type Interpreter interface {
	Interpret(
		name string,
		value interface{},
		param *model.Param,
		parentOpHandle model.DataHandle,
		scope map[string]*model.Value,
		opScratchDir string,
	) (*model.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		arrayInterpreter:   array.NewInterpreter(),
		booleanInterpreter: boolean.NewInterpreter(),
		dirInterpreter:     dir.NewInterpreter(),
		fileInterpreter:    file.NewInterpreter(),
		numberInterpreter:  number.NewInterpreter(),
		objectInterpreter:  object.NewInterpreter(),
		stringInterpreter:  stringPkg.NewInterpreter(),
	}
}

type _interpreter struct {
	arrayInterpreter   array.Interpreter
	booleanInterpreter boolean.Interpreter
	dirInterpreter     dir.Interpreter
	fileInterpreter    file.Interpreter
	numberInterpreter  number.Interpreter
	objectInterpreter  object.Interpreter
	stringInterpreter  stringPkg.Interpreter
}

func (itp _interpreter) Interpret(
	name string,
	valueExpression interface{},
	param *model.Param,
	parentOpHandle model.DataHandle,
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
		arrayValue, err := itp.arrayInterpreter.Interpret(scope, valueExpression, parentOpHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return arrayValue, nil
	case nil != param.Boolean:
		booleanValue, err := itp.booleanInterpreter.Interpret(scope, valueExpression, parentOpHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return booleanValue, nil
	case nil != param.File:
		fileValue, err := itp.fileInterpreter.Interpret(scope, valueExpression, parentOpHandle, opScratchDir)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return fileValue, nil
	case nil != param.Dir:
		dirValue, err := itp.dirInterpreter.Interpret(scope, valueExpression, parentOpHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return dirValue, nil
	case nil != param.Number:
		numberValue, err := itp.numberInterpreter.Interpret(scope, valueExpression, parentOpHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return numberValue, nil
	case nil != param.Object:
		objectValue, err := itp.objectInterpreter.Interpret(scope, valueExpression, parentOpHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return objectValue, nil
	case nil != param.String:
		stringValue, err := itp.stringInterpreter.Interpret(scope, valueExpression, parentOpHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		return stringValue, nil
	case nil != param.Socket:
		return nil, fmt.Errorf("unable to bind '%v' to '%+v'; sockets must be passed by reference", name, valueExpression)
	}

	return nil, fmt.Errorf("unable to bind '%v' to '%v'", name, valueExpression)
}
