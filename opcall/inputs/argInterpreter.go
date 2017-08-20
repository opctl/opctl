package inputs

//go:generate counterfeiter -o ./fakeArgInterpreter.go --fake-name fakeArgInterpreter ./ argInterpreter

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	stringPkg "github.com/opspec-io/sdk-golang/string"
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
		parentPkgRef string,
		scope map[string]*model.Value,
	) (*model.Value, error)
}

func newArgInterpreter() argInterpreter {
	return _argInterpreter{
		interpolater: interpolater.New(),
		string:       stringPkg.New(),
	}
}

type _argInterpreter struct {
	interpolater interpolater.Interpolater
	string       stringPkg.String
}

func (ai _argInterpreter) Interpret(
	name string,
	value interface{},
	param *model.Param,
	parentPkgRef string,
	scope map[string]*model.Value,
) (*model.Value, error) {

	if nil == param {
		return nil, fmt.Errorf("Unable to bind to '%v'; '%v' not a defined input", name, name)
	}

	if nil == value || "" == value {
		// implicitly bound
		dcgValue, ok := scope[name]
		if !ok {
			return nil, fmt.Errorf("Unable to bind to '%v' via implicit ref; '%v' not in scope", name, name)
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

	if stringExpression, ok := value.(string); ok {
		if deprecatedExplicitRef, ok := scope[stringExpression]; ok {
			// deprecated explicit arg
			return deprecatedExplicitRef, nil
		} else if explicitRef := strings.TrimSuffix(strings.TrimPrefix(stringExpression, "$("), ")"); len(explicitRef) == (len(stringExpression) - 3) {
			// @FIXME this won't handle $(ref1)$(ref2) properly
			// explicit arg
			dcgValue, ok := scope[explicitRef]
			if !ok {
				return nil, fmt.Errorf("Unable to bind '%v' to '%v' via explicit ref; '%v' not in scope", name, explicitRef, explicitRef)
			}
			return dcgValue, nil
		} else {
			interpolatedVal := ai.interpolater.Interpolate(stringExpression, scope)
			switch {
			// interpolated arg
			case nil != param.String:
				stringValue, err := ai.string.Interpret(scope, stringExpression)
				if nil != err {
					return nil, fmt.Errorf("Unable to bind '%v' to '%v'; error was: '%v'", name, stringValue, err.Error())
				}
				return &model.Value{String: &stringValue}, nil
			case nil != param.Dir:
				if strings.HasPrefix(stringExpression, "/") {
					// bound to pkg dir
					interpolatedVal = filepath.Join(parentPkgRef, interpolatedVal)
				}
				return &model.Value{Dir: ai.rootPath(interpolatedVal)}, nil
			case nil != param.Number:
				floatVal, err := strconv.ParseFloat(interpolatedVal, 64)
				if nil != err {
					return nil, fmt.Errorf("Unable to bind '%v' to '%v' as number; error was: '%v'", name, interpolatedVal, err.Error())
				}
				return &model.Value{Number: &floatVal}, nil
			case nil != param.File:
				if strings.HasPrefix(stringExpression, "/") {
					// bound to pkg file
					interpolatedVal = filepath.Join(parentPkgRef, interpolatedVal)
				}
				return &model.Value{File: ai.rootPath(interpolatedVal)}, nil
			case nil != param.Socket:
				return nil, fmt.Errorf("Unable to bind '%v' to '%v'; sockets must be passed by reference", name, interpolatedVal)
			}
		}
	}

	return nil, fmt.Errorf("Unable to bind '%v' to '%v'", name, value)
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
