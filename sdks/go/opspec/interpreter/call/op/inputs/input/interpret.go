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

	if param == nil {
		return nil, fmt.Errorf("unable to bind to '%v': '%v' not a defined input", name, name)
	}

	if valueExpression == nil || valueExpression == "" {
		// implicitly bound
		_, ok := scope[name]
		if !ok {
			return nil, fmt.Errorf("unable to bind to '%v' via implicit ref: '%v' not in scope", name, name)
		}
		valueExpression = opspec.NameToRef(name)
	}

	switch {
	case param.Array != nil:
		arrayValue, err := array.Interpret(scope, valueExpression)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return arrayValue, nil
	case param.Boolean != nil:
		booleanValue, err := boolean.Interpret(scope, valueExpression)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return booleanValue, nil
	case param.File != nil:
		fileValue, err := file.Interpret(scope, valueExpression, opScratchDir, true)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return fileValue, nil
	case param.Dir != nil:
		dirValue, err := dir.Interpret(scope, valueExpression, opScratchDir, true)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return dirValue, nil
	case param.Number != nil:
		numberValue, err := number.Interpret(scope, valueExpression)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return numberValue, nil
	case param.Object != nil:
		objectValue, err := object.Interpret(scope, valueExpression)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return objectValue, nil
	case param.String != nil:
		stringValue, err := str.Interpret(scope, valueExpression)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		return stringValue, nil
	case param.Socket != nil:
		stringValueExpression, isString := valueExpression.(string)
		if !isString {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v': sockets must be passed by reference", name, valueExpression)
		}

		socketValue, err := reference.Interpret(stringValueExpression, scope, nil)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to bind '%v' to '%+v'", name, valueExpression))
		}
		if socketValue.Socket == nil {
			return nil, fmt.Errorf("unable to bind '%v' to '%+v': '%+v' must reference a socket", name, valueExpression, valueExpression)
		}

		return socketValue, nil
	}

	return nil, fmt.Errorf("unable to bind '%v' to '%v'", name, valueExpression)
}
