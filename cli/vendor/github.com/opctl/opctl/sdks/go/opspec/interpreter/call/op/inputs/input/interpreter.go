package input

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/array"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/boolean"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/dir"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/file"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/number"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/object"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
	stringPkg "github.com/opctl/opctl/sdks/go/opspec/interpreter/string"
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
		arrayInterpreter:     array.NewInterpreter(),
		booleanInterpreter:   boolean.NewInterpreter(),
		dirInterpreter:       dir.NewInterpreter(),
		fileInterpreter:      file.NewInterpreter(),
		numberInterpreter:    number.NewInterpreter(),
		objectInterpreter:    object.NewInterpreter(),
		referenceInterpreter: reference.NewInterpreter(),
		stringInterpreter:    stringPkg.NewInterpreter(),
	}
}

type _interpreter struct {
	arrayInterpreter     array.Interpreter
	booleanInterpreter   boolean.Interpreter
	dirInterpreter       dir.Interpreter
	fileInterpreter      file.Interpreter
	numberInterpreter    number.Interpreter
	objectInterpreter    object.Interpreter
	referenceInterpreter reference.Interpreter
	stringInterpreter    stringPkg.Interpreter
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
		_, ok := scope[name]
		if !ok {
			return nil, fmt.Errorf("unable to bind to '%v' via implicit ref; '%v' not in scope", name, name)
		}
		valueExpression = fmt.Sprintf("$(%v)", name)
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
		stringValueExpression, isString := valueExpression.(string)
		if !isString {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; sockets must be passed by reference", name, valueExpression)
		}

		socketValue, err := itp.referenceInterpreter.Interpret(stringValueExpression, scope, parentOpHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; error was: '%v'", name, valueExpression, err.Error())
		}
		if nil == socketValue.Socket {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v'; '%+v' must reference a socket", name, valueExpression, valueExpression)
		}

		return socketValue, nil
	}

	return nil, fmt.Errorf("unable to bind '%v' to '%v'", name, valueExpression)
}
