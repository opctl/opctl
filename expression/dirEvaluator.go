package expression

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"strings"
)

type dirEvaluator interface {
	// EvalToDir evaluates an expression to a dir value
	//
	// Examples of valid dir expressions:
	// scope ref: $(scope-ref)
	// scope ref w/ path: $(scope-ref/sub-dir)
	// scope ref w/ deprecated path: $(scope-ref)/sub-dir
	// pkg fs ref: $(/pkg-fs-ref)
	// pkg fs ref w/ path: $(/pkg-fs-ref/sub-dir)
	EvalToDir(
		scope map[string]*model.Value,
		expression string,
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func newDirEvaluator() dirEvaluator {
	return _dirEvaluator{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _dirEvaluator struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (etd _dirEvaluator) EvalToDir(
	scope map[string]*model.Value,
	expression string,
	pkgHandle model.PkgHandle,
) (*model.Value, error) {
	possibleRefCloserIndex := strings.Index(expression, interpolater.RefEnd)

	// the following is gross but it's due to all the deprecated syntax we need to handle
	if strings.HasPrefix(expression, "/") {

		// deprecated pkg fs ref
		pkgFsRefPath, err := etd.interpolater.Interpolate(
			expression,
			scope,
			pkgHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("unable to evaluate %v to dir; error was %v", expression, err.Error())
		}

		pkgFsRefPath = filepath.Join(pkgHandle.Ref(), pkgFsRefPath)
		return &model.Value{Dir: &pkgFsRefPath}, err

	} else if strings.HasPrefix(expression, interpolater.RefStart) && possibleRefCloserIndex > 0 {

		refExpression := expression[2:possibleRefCloserIndex]
		refParts := strings.SplitN(refExpression, "/", 2)
		var dirValue string

		if strings.HasPrefix(refExpression, "/") {

			// pkg fs ref
			pkgFsRef, err := etd.interpolater.Interpolate(refExpression, scope, pkgHandle)
			if nil != err {
				return nil, fmt.Errorf("unable to evaluate pkg fs ref %v; error was %v", refExpression, err.Error())
			}
			dirValue = filepath.Join(pkgHandle.Ref(), pkgFsRef)

		} else if dcgValue, ok := scope[refExpression]; ok && nil != dcgValue.Dir {

			// scope ref
			dirValue = *dcgValue.Dir

		} else if dcgValue, ok := scope[refParts[0]]; ok && nil != dcgValue.Dir {

			// scope ref w/ path
			pathExpression := refParts[1]
			path, err := etd.interpolater.Interpolate(pathExpression, scope, pkgHandle)
			if nil != err {
				return nil, fmt.Errorf("unable to evaluate path %v; error was %v", pathExpression, err.Error())
			}
			dirValue = filepath.Join(*dcgValue.Dir, path)

		}

		if len(expression) > possibleRefCloserIndex+1 {
			// evaluate deprecated path
			deprecatedPathExpression := expression[possibleRefCloserIndex+1:]
			deprecatedPath, err := etd.interpolater.Interpolate(deprecatedPathExpression, scope, pkgHandle)
			if nil != err {
				return nil, fmt.Errorf("unable to evaluate path %v; error was %v", deprecatedPathExpression, err.Error())
			}

			dirValue := filepath.Join(dirValue, deprecatedPath)
			return &model.Value{Dir: &dirValue}, nil
		}

		return &model.Value{Dir: &dirValue}, nil
	}

	return nil, fmt.Errorf("unable to evaluate %v to dir", expression)

}
