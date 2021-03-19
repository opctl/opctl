package input

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/array"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/boolean"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/dir"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/file"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/number"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/object"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
	"github.com/pkg/errors"
)

//Interpret interprets an op input
func Interpret(
	name string,
	valueExpression interface{},
	param *model.Param,
	scope map[string]*model.Value,
	opScratchDir string,
) (*model.Value, error) {

	if nil == param {
		return nil, fmt.Errorf("unable to bind to '%v': '%v' not a defined input", name, name)
	}

	if nil == valueExpression || "" == valueExpression {
		// implicitly bound
		_, ok := scope[name]
		if !ok {
			return nil, fmt.Errorf("unable to bind to '%v' via implicit ref: '%v' not in scope", name, name)
		}
		valueExpression = opspec.NameToRef(name)
	}

	switch {
	case nil != param.Array:
		arrayValue, err := array.Interpret(scope, valueExpression)
		if nil != err {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return arrayValue, nil
	case nil != param.Boolean:
		booleanValue, err := boolean.Interpret(scope, valueExpression)
		if nil != err {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return booleanValue, nil
	case nil != param.File:
		fileValue, err := file.Interpret(scope, valueExpression, opScratchDir, true)
		if nil != err {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return fileValue, nil
	case nil != param.Dir:
		dirValue, err := dir.Interpret(scope, valueExpression, opScratchDir, true)
		if nil != err {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return dirValue, nil
	case nil != param.Number:
		numberValue, err := number.Interpret(scope, valueExpression)
		if nil != err {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return numberValue, nil
	case nil != param.Object:
		objectValue, err := object.Interpret(scope, valueExpression)
		if nil != err {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return objectValue, nil
	case nil != param.String:
		stringValue, err := str.Interpret(scope, valueExpression)
		if nil != err {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return stringValue, nil
	case nil != param.Socket:
		stringValueExpression, isString := valueExpression.(string)
		if !isString {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v': sockets must be passed by reference", name, valueExpression)
		}

		socketValue, err := reference.Interpret(stringValueExpression, scope, nil)
		if nil != err {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		if nil == socketValue.Socket {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v': '%+v' must reference a socket", name, valueExpression, valueExpression)
		}

		return socketValue, nil
	}

	return nil, fmt.Errorf("unable to bind '%v' to '%v'", name, valueExpression)
}
