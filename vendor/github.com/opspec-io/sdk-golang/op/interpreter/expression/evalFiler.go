package expression

import (
	"fmt"
	"github.com/golang-interfaces/iio"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression/interpolater"
	"path/filepath"
	"strings"
)

type evalFiler interface {
	// EvalToFile evaluates an expression to a file value
	// expression must be a type supported by coerce.ToFile
	// scratchDir will be used as the containing dir if file creation necessary
	//
	// Examples of valid file expressions:
	// scope ref: $(scope-ref)
	// scope ref w/ path: $(scope-ref/file.txt)
	// scope ref w/ deprecated path: $(scope-ref)/file.txt
	// pkg fs ref: $(/pkg-fs-ref)
	// pkg fs ref w/ path: $(/pkg-fs-ref/file.txt)
	EvalToFile(
		scope map[string]*model.Value,
		expression interface{},
		opDirHandle model.DataHandle,
		scratchDir string,
	) (*model.Value, error)
}

func newEvalFiler() evalFiler {
	return _evalFiler{
		evalArrayInitializerer:  newEvalArrayInitializerer(),
		evalObjectInitializerer: newEvalObjectInitializerer(),
		coerce:                  coerce.New(),
		interpolater:            interpolater.New(),
		io:                      iio.New(),
		os:                      ios.New(),
	}
}

type _evalFiler struct {
	evalArrayInitializerer
	evalObjectInitializerer
	coerce       coerce.Coerce
	interpolater interpolater.Interpolater
	io           iio.IIO
	os           ios.IOS
}

func (ef _evalFiler) EvalToFile(
	scope map[string]*model.Value,
	expression interface{},
	opDirHandle model.DataHandle,
	scratchDir string,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case float64:
		return ef.coerce.ToFile(&model.Value{Number: &expression}, scratchDir)
	case map[string]interface{}:
		objectValue, err := ef.evalObjectInitializerer.Eval(
			expression,
			scope,
			opDirHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %+v to file; error was %v", expression, err)
		}

		return ef.coerce.ToFile(&model.Value{Object: objectValue}, scratchDir)
	case []interface{}:
		arrayValue, err := ef.evalArrayInitializerer.Eval(
			expression,
			scope,
			opDirHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %+v to file; error was %v", expression, err)
		}

		return ef.coerce.ToFile(&model.Value{Array: arrayValue}, scratchDir)
	case string:

		possibleRefCloserIndex := strings.Index(expression, interpolater.RefEnd)
		if ref, ok := tryResolveExplicitRef(expression, scope); ok {
			// scope ref w/out path
			return ef.coerce.ToFile(ref, scratchDir)
		} else if strings.HasPrefix(expression, interpolater.RefStart) && possibleRefCloserIndex > 0 {

			refExpression := expression[2:possibleRefCloserIndex]
			refParts := strings.SplitN(refExpression, "/", 2)

			if strings.HasPrefix(refExpression, "/") && len(expression) == possibleRefCloserIndex+1 {

				// pkg fs ref
				pkgFsRef, err := ef.interpolater.Interpolate(refExpression, scope, opDirHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to evaluate pkg fs ref %v; error was %v", refExpression, err.Error())
				}

				fileValue := filepath.Join(*opDirHandle.Path(), pkgFsRef)

				return &model.Value{File: &fileValue}, nil

			} else if dcgValue, ok := scope[refParts[0]]; ok && nil != dcgValue.Dir {

				// dir scope ref w/ path
				pathExpression := refParts[1]
				path, err := ef.interpolater.Interpolate(pathExpression, scope, opDirHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to evaluate path %v; error was %v", pathExpression, err.Error())
				}

				fileValue := filepath.Join(*dcgValue.Dir, path)
				return &model.Value{File: &fileValue}, nil

			}

		}
		// plain string
		stringValue, err := ef.interpolater.Interpolate(expression, scope, opDirHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %v to file; error was %v", expression, err.Error())
		}
		return ef.coerce.ToFile(&model.Value{String: &stringValue}, scratchDir)
	}

	return nil, fmt.Errorf("unable to evaluate %+v to file", expression)

}
