package file

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/opctl/sdk-golang/opspec/interpreter"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"

	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/value"
)

type Interpreter interface {
	// Interpret interprets an expression to a file value.
	// Expression must be a type supported by coerce.ToFile
	// scratchDir will be used as the containing dir if file creation necessary
	//
	// Examples of valid file expressions:
	// scope ref: $(scope-ref)
	// scope ref w/ path: $(scope-ref/file.txt)
	// pkg fs ref: $(/pkg-fs-ref)
	// pkg fs ref w/ path: $(/pkg-fs-ref/file.txt)
	Interpret(
		scope map[string]*model.Value,
		expression interface{},
		opHandle model.DataHandle,
		scratchDir string,
	) (*model.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		coerce:           coerce.New(),
		interpolater:     interpolater.New(),
		valueInterpreter: value.NewInterpreter(),
	}
}

type _interpreter struct {
	coerce           coerce.Coerce
	interpolater     interpolater.Interpolater
	valueInterpreter value.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	expression interface{},
	opHandle model.DataHandle,
	scratchDir string,
) (*model.Value, error) {
	if expressionAsString, ok := expression.(string); ok {
		if ref, ok := interpreter.TryResolveExplicitRef(expressionAsString, scope); ok {
			// scope ref w/out path
			return itp.coerce.ToFile(ref, scratchDir)
		} else if strings.HasPrefix(expressionAsString, interpolater.RefStart) && strings.HasSuffix(expressionAsString, interpolater.RefEnd) {

			refExpression := strings.TrimSuffix(strings.TrimPrefix(expressionAsString, interpolater.RefStart), interpolater.RefEnd)
			refParts := strings.SplitN(refExpression, "/", 2)

			if strings.HasPrefix(refExpression, "/") {

				// pkg fs ref
				pkgFsRef, err := itp.interpolater.Interpolate(refExpression, scope, opHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to evaluate pkg fs ref %v; error was %v", refExpression, err.Error())
				}

				fileValue := filepath.Join(*opHandle.Path(), pkgFsRef)

				return &model.Value{File: &fileValue}, nil

			} else if dcgValue, ok := scope[refParts[0]]; ok && nil != dcgValue.Dir {

				// dir scope ref w/ path
				pathExpression := refParts[1]
				path, err := itp.interpolater.Interpolate(pathExpression, scope, opHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to evaluate path %v; error was %v", pathExpression, err.Error())
				}

				fileValue := filepath.Join(*dcgValue.Dir, path)
				return &model.Value{File: &fileValue}, nil

			}

		}
	}

	value, err := itp.valueInterpreter.Interpret(
		expression,
		scope,
		opHandle,
	)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret %+v to file; error was %v", expression, err)
	}

	return itp.coerce.ToFile(&value, scratchDir)
}
