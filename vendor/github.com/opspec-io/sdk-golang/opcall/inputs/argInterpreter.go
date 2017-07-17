package inputs

//go:generate counterfeiter -o ./fakeArgInterpreter.go --fake-name fakeArgInterpreter ./ argInterpreter

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type argInterpreter interface {
	Interpret(
		name,
		value string,
		param *model.Param,
		pkgPath string,
		scope map[string]*model.Value,
	) (*model.Value, error)
}

func newArgInterpreter() argInterpreter {
	return _argInterpreter{
		interpolater: interpolater.New(),
	}
}

type _argInterpreter struct {
	interpolater interpolater.Interpolater
}

func (ai _argInterpreter) Interpret(
	name,
	value string,
	param *model.Param,
	pkgPath string,
	scope map[string]*model.Value,
) (*model.Value, error) {

	var dcgValue *model.Value
	if nil == param {
		return nil, fmt.Errorf("Unable to bind to '%v'; '%v' not a defined input", name, name)
	} else if "" == value {
		// implicit arg
		var ok bool
		dcgValue, ok = scope[name]
		if !ok {
			return nil, fmt.Errorf("Unable to bind to '%v' via implicit ref; '%v' not in scope", name, name)
		}
	} else if deprecatedExplicitRef, ok := scope[value]; ok {
		// deprecated explicit arg
		dcgValue = deprecatedExplicitRef
	} else if explicitRef := strings.TrimSuffix(strings.TrimPrefix(value, "$("), ")"); len(explicitRef) == (len(value) - 3) {
		// explicit arg
		dcgValue, ok = scope[explicitRef]
		if !ok {
			return nil, fmt.Errorf("Unable to bind '%v' to '%v' via explicit ref; '%v' not in scope", name, explicitRef, explicitRef)
		}
	} else {
		interpolatedVal := ai.interpolater.Interpolate(value, scope)
		switch {
		// interpolated arg
		case nil != param.String:
			dcgValue = &model.Value{String: &interpolatedVal}
		case nil != param.Dir:
			if strings.HasPrefix(value, "/") {
				// bound to pkg dir
				interpolatedVal = filepath.Join(pkgPath, interpolatedVal)
			}
			dcgValue = &model.Value{Dir: ai.rootPath(interpolatedVal)}
		case nil != param.Number:
			floatVal, err := strconv.ParseFloat(interpolatedVal, 64)
			if nil != err {
				return nil, fmt.Errorf("Unable to bind '%v' to '%v' as number; error was: '%v'", name, interpolatedVal, err.Error())
			}
			dcgValue = &model.Value{Number: &floatVal}
		case nil != param.File:
			if strings.HasPrefix(value, "/") {
				// bound to pkg file
				interpolatedVal = filepath.Join(pkgPath, interpolatedVal)
			}
			dcgValue = &model.Value{File: ai.rootPath(interpolatedVal)}
		case nil != param.Socket:
			return nil, fmt.Errorf("Unable to bind '%v' to '%v'; sockets must be passed by reference", name, interpolatedVal)
		}
	}

	return dcgValue, nil
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
