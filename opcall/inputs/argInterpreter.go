package inputs

//go:generate counterfeiter -o ./fakeArgInterpreter.go --fake-name fakeArgInterpreter ./ argInterpreter

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/expression"
	"github.com/opspec-io/sdk-golang/model"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type argInterpreter interface {
	Interpret(
		name string,
		value interface{},
		param *model.Param,
		parentPkgHandle model.PkgHandle,
		scope map[string]*model.Value,
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
			// deprecated explicit arg
			return dcgValue, nil
		} else if dcgValue, ok := scope[strings.TrimSuffix(strings.TrimPrefix(stringValue, "$("), ")")]; ok {
			// explicit arg
			return dcgValue, nil
		} else {
			switch {
			// interpolated arg
			case nil != param.String:
				stringValue, err := ai.expression.EvalToString(scope, stringValue, parentPkgHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to bind '%v' to '%v'; error was: '%v'", name, stringValue, err.Error())
				}
				return &model.Value{String: &stringValue}, nil
			case nil != param.Dir:
				interpolatedVal, err := ai.expression.EvalToString(scope, stringValue, parentPkgHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to bind '%v' to '%v'; error was: '%v'", name, stringValue, err.Error())
				}
				return &model.Value{Dir: ai.rootPath(interpolatedVal)}, nil
			case nil != param.Number:
				numberValue, err := ai.expression.EvalToNumber(scope, stringValue, parentPkgHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to bind '%v' to '%v'; error was: '%v'", name, stringValue, err.Error())
				}
				return &model.Value{Number: &numberValue}, nil
			case nil != param.File:
				interpolatedVal, err := ai.expression.EvalToString(scope, stringValue, parentPkgHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to bind '%v' to '%v'; error was: '%v'", name, stringValue, err.Error())
				}
				return &model.Value{File: ai.rootPath(interpolatedVal)}, nil
			case nil != param.Socket:
				return nil, fmt.Errorf("unable to bind '%v' to '%v'; sockets must be passed by reference", name, stringValue)
			}
		}
	}

	return nil, fmt.Errorf("unable to bind '%v' to '%v'", name, value)
}

// rootPath ensures paths are rooted (interpreted as having no parent) so parent paths of input files/dirs aren't
// accessible (which would break encapsulation)
func (ai _argInterpreter) rootPath(
	path string,
) *string {
	path = strings.Replace(path, "../", string(os.PathSeparator), -1)
	path = strings.Replace(path, "..\\", string(os.PathSeparator), -1)
	path = filepath.Clean(path)
	return &path
}
