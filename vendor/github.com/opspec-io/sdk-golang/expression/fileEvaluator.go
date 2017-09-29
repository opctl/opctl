package expression

import (
	"fmt"
	"github.com/golang-interfaces/iio"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"path/filepath"
	"strings"
)

type fileEvaluator interface {
	// EvalToFile evaluates an expression to a file value
	// expression must be a type supported by data.CoerceToFile
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
		pkgHandle model.PkgHandle,
		scratchDir string,
	) (*model.Value, error)
}

func newFileEvaluator() fileEvaluator {
	return _fileEvaluator{
		data:         data.New(),
		interpolater: interpolater.New(),
		io:           iio.New(),
		os:           ios.New(),
		pkg:          pkg.New(),
	}
}

type _fileEvaluator struct {
	data         data.Data
	interpolater interpolater.Interpolater
	io           iio.IIO
	os           ios.IOS
	pkg          pkg.Pkg
}

func (etf _fileEvaluator) EvalToFile(
	scope map[string]*model.Value,
	expression interface{},
	pkgHandle model.PkgHandle,
	scratchDir string,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case float64:
		return etf.data.CoerceToFile(&model.Value{Number: &expression}, scratchDir)
	case map[string]interface{}:
		return etf.data.CoerceToFile(&model.Value{Object: expression}, scratchDir)
	case []interface{}:
		return etf.data.CoerceToFile(&model.Value{Array: expression}, scratchDir)
	case string:

		// this block is gross but it's due to the deprecated syntax we support
		possibleRefCloserIndex := strings.Index(expression, interpolater.RefEnd)
		if ref, ok := tryResolveExplicitRef(expression, scope); ok {
			// scope ref w/out path
			return etf.data.CoerceToFile(ref, scratchDir)
		} else if strings.HasPrefix(expression, "/") {
			// deprecated pkg fs ref
			deprecatedPkgFsRefPath, err := etf.interpolater.Interpolate(
				expression,
				scope,
				pkgHandle,
			)
			if nil != err {
				return nil, fmt.Errorf("unable to evaluate %v to file; error was %v", expression, err.Error())
			}

			fileValue := filepath.Join(*pkgHandle.Path(), deprecatedPkgFsRefPath)

			return &model.Value{File: &fileValue}, nil

		} else if strings.HasPrefix(expression, interpolater.RefStart) && possibleRefCloserIndex > 0 {

			refExpression := expression[2:possibleRefCloserIndex]
			refParts := strings.SplitN(refExpression, "/", 2)

			if strings.HasPrefix(refExpression, "/") && len(expression) == possibleRefCloserIndex+1 {

				// pkg fs ref
				pkgFsRef, err := etf.interpolater.Interpolate(refExpression, scope, pkgHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to evaluate pkg fs ref %v; error was %v", refExpression, err.Error())
				}

				fileValue := filepath.Join(*pkgHandle.Path(), pkgFsRef)

				return &model.Value{File: &fileValue}, nil

			} else if dcgValue, ok := scope[refExpression]; ok && nil != dcgValue.Dir {

				// dir scope ref w/ deprecated path
				deprecatedPathExpression := expression[possibleRefCloserIndex+1:]
				deprecatedPath, err := etf.interpolater.Interpolate(deprecatedPathExpression, scope, pkgHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to evaluate path %v; error was %v", deprecatedPathExpression, err.Error())
				}

				fileValue := filepath.Join(*dcgValue.Dir, deprecatedPath)
				return &model.Value{File: &fileValue}, nil

			} else if dcgValue, ok := scope[refParts[0]]; ok && nil != dcgValue.Dir {

				// dir scope ref w/ path
				pathExpression := refParts[1]
				path, err := etf.interpolater.Interpolate(pathExpression, scope, pkgHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to evaluate path %v; error was %v", pathExpression, err.Error())
				}

				// no possibility of deprecated path due to presence of path
				fileValue := filepath.Join(*dcgValue.Dir, path)
				return &model.Value{File: &fileValue}, nil

			}
			// non-dir ref suffixed by ref(s) &/or non-refs

		}
		// plain string
		stringValue, err := etf.interpolater.Interpolate(expression, scope, pkgHandle)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %v to file; error was %v", expression, err.Error())
		}
		return etf.data.CoerceToFile(&model.Value{String: &stringValue}, scratchDir)
	}

	return nil, fmt.Errorf("unable to evaluate %+v to file", expression)

}
