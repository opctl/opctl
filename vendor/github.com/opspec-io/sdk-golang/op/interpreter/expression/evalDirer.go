package expression

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression/interpolater"
)

type evalDirer interface {
	// EvalToDir evaluates an expression to a dir value
	// expression must be of type string
	//
	// Examples of valid dir expressions:
	// scope ref: $(scope-ref)
	// scope ref w/ path: $(scope-ref/sub-dir)
	// pkg fs ref: $(/pkg-fs-ref)
	// pkg fs ref w/ path: $(/pkg-fs-ref/sub-dir)
	EvalToDir(
		scope map[string]*model.Value,
		expression interface{},
		opDirHandle model.DataHandle,
	) (*model.Value, error)
}

func newEvalDirer() evalDirer {
	return _evalDirer{
		coerce:       coerce.New(),
		interpolater: interpolater.New(),
	}
}

type _evalDirer struct {
	coerce       coerce.Coerce
	interpolater interpolater.Interpolater
}

func (ed _evalDirer) EvalToDir(
	scope map[string]*model.Value,
	expression interface{},
	opDirHandle model.DataHandle,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case string:
		possibleRefCloserIndex := strings.Index(expression, interpolater.RefEnd)

		if ref, ok := tryResolveExplicitRef(expression, scope); ok && nil != ref.Dir {
			// scope ref w/out path
			return ref, nil
		} else if strings.HasPrefix(expression, interpolater.RefStart) && possibleRefCloserIndex > 0 {

			refExpression := expression[2:possibleRefCloserIndex]
			refParts := strings.SplitN(refExpression, "/", 2)
			var dirValue string

			if strings.HasPrefix(refExpression, "/") {

				// pkg fs ref
				pkgFsRef, err := ed.interpolater.Interpolate(refExpression, scope, opDirHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to evaluate pkg fs ref %v; error was %v", refExpression, err.Error())
				}
				pkgPath := opDirHandle.Path()
				dirValue = filepath.Join(*pkgPath, pkgFsRef)

			} else if dcgValue, ok := scope[refParts[0]]; ok && nil != dcgValue.Dir {

				// scope ref w/ path
				pathExpression := refParts[1]
				path, err := ed.interpolater.Interpolate(pathExpression, scope, opDirHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to evaluate path %v; error was %v", pathExpression, err.Error())
				}

				dirValue := filepath.Join(*dcgValue.Dir, path)
				return &model.Value{Dir: &dirValue}, nil

			}

			return &model.Value{Dir: &dirValue}, nil
		}
	}

	return nil, fmt.Errorf("unable to evaluate %v to dir", expression)

}
