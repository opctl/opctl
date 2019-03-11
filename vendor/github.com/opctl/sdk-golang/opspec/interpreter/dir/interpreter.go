package dir

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"github.com/opctl/sdk-golang/opspec/interpreter"
	"path/filepath"
	"strings"

	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"
)

type Interpreter interface {
	// Interpret evaluates an expression to a dir value.
	// Expression must be of type string.
	//
	// Examples of valid dir expressions:
	// scope ref: $(scope-ref)
	// scope ref w/ path: $(scope-ref/sub-dir)
	// pkg fs ref: $(/pkg-fs-ref)
	// pkg fs ref w/ path: $(/pkg-fs-ref/sub-dir)
	Interpret(
		scope map[string]*model.Value,
		expression interface{},
		opHandle model.DataHandle,
	) (*model.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		coerce:       coerce.New(),
		interpolater: interpolater.New(),
	}
}

type _interpreter struct {
	coerce       coerce.Coerce
	interpolater interpolater.Interpolater
}

func (ed _interpreter) Interpret(
	scope map[string]*model.Value,
	expression interface{},
	opHandle model.DataHandle,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case string:
		possibleRefCloserIndex := strings.Index(expression, interpolater.RefEnd)

		if ref, ok := interpreter.TryResolveExplicitRef(expression, scope); ok && nil != ref.Dir {
			// scope ref w/out path
			return ref, nil
		} else if strings.HasPrefix(expression, interpolater.RefStart) && possibleRefCloserIndex > 0 {

			refExpression := expression[2:possibleRefCloserIndex]
			refParts := strings.SplitN(refExpression, "/", 2)
			var dirValue string

			if strings.HasPrefix(refExpression, "/") {

				// pkg fs ref
				pkgFsRef, err := ed.interpolater.Interpolate(refExpression, scope, opHandle)
				if nil != err {
					return nil, fmt.Errorf("unable to evaluate pkg fs ref %v; error was %v", refExpression, err.Error())
				}
				opPath := opHandle.Path()
				dirValue = filepath.Join(*opPath, pkgFsRef)

			} else if dcgValue, ok := scope[refParts[0]]; ok && nil != dcgValue.Dir {

				// scope ref w/ path
				pathExpression := refParts[1]
				path, err := ed.interpolater.Interpolate(pathExpression, scope, opHandle)
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
