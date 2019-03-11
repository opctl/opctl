package file

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"github.com/golang-interfaces/iio"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter"
	arrayInitializer "github.com/opctl/sdk-golang/opspec/interpreter/array/initializer"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"
	objectInitializer "github.com/opctl/sdk-golang/opspec/interpreter/object/initializer"
	"path/filepath"
	"strings"
)

type Interpreter interface {
	// Interpret evaluates an expression to a file value.
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
		arrayInitializerInterpreter:  arrayInitializer.NewInterpreter(),
		objectInitializerInterpreter: objectInitializer.NewInterpreter(),
		coerce:       coerce.New(),
		interpolater: interpolater.New(),
		io:           iio.New(),
		os:           ios.New(),
	}
}

type _interpreter struct {
	arrayInitializerInterpreter  arrayInitializer.Interpreter
	objectInitializerInterpreter objectInitializer.Interpreter
	coerce                       coerce.Coerce
	interpolater                 interpolater.Interpolater
	io                           iio.IIO
	os                           ios.IOS
}

func (ef _interpreter) Interpret(
	scope map[string]*model.Value,
	expression interface{},
	opHandle model.DataHandle,
	scratchDir string,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case float64:
		return ef.coerce.ToFile(&model.Value{Number: &expression}, scratchDir)
	case int:
		expressionAsFloat64 := float64(expression)
		return ef.coerce.ToFile(&model.Value{Number: &expressionAsFloat64}, scratchDir)
	case map[string]interface{}:
		objectValue, err := ef.objectInitializerInterpreter.Interpret(
			expression,
			scope,
			opHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %+v to file; error was %v", expression, err)
		}

		return ef.coerce.ToFile(&model.Value{Object: objectValue}, scratchDir)
	case []interface{}:
		arrayValue, err := ef.arrayInitializerInterpreter.Interpret(
			expression,
			scope,
			opHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %+v to file; error was %v", expression, err)
		}

		return ef.coerce.ToFile(&model.Value{Array: arrayValue}, scratchDir)
	case string:

		possibleRefCloserIndex := strings.Index(expression, interpolater.RefEnd)
		if ref, ok := interpreter.TryResolveExplicitRef(expression, scope); ok {
			// scope ref w/out path
			return ef.coerce.ToFile(ref, scratchDir)
		} else if strings.HasPrefix(expression, interpolater.RefStart) && possibleRefCloserIndex > 0 {

			refExpression := expression[2:possibleRefCloserIndex]
			refParts := strings.SplitN(refExpression, "/", 2)

			if strings.HasPrefix(refExpression, "/") && len(expression) == possibleRefCloserIndex+1 {

				// pkg fs ref
				pkgFsRef, err := ef.interpolater.Interpolate(refExpression, scope, opHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to evaluate pkg fs ref %v; error was %v", refExpression, err.Error())
				}

				fileValue := filepath.Join(*opHandle.Path(), pkgFsRef)

				return &model.Value{File: &fileValue}, nil

			} else if dcgValue, ok := scope[refParts[0]]; ok && nil != dcgValue.Dir {

				// dir scope ref w/ path
				pathExpression := refParts[1]
				path, err := ef.interpolater.Interpolate(pathExpression, scope, opHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to evaluate path %v; error was %v", pathExpression, err.Error())
				}

				fileValue := filepath.Join(*dcgValue.Dir, path)
				return &model.Value{File: &fileValue}, nil

			}

		}
		// plain string
		stringValue, err := ef.interpolater.Interpolate(expression, scope, opHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %v to file; error was %v", expression, err.Error())
		}
		return ef.coerce.ToFile(&model.Value{String: &stringValue}, scratchDir)
	}

	return nil, fmt.Errorf("unable to evaluate %+v to file", expression)

}
