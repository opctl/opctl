package dir

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"strings"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
	"github.com/opctl/opctl/sdks/go/types"
)

type Interpreter interface {
	// Interpret interprets an expression to a dir value.
	// Expression must be of type string.
	//
	// Examples of valid dir expressions:
	// scope ref: $(scope-ref)
	// scope ref w/ path: $(scope-ref/sub-dir)
	// pkg fs ref: $(/pkg-fs-ref)
	// pkg fs ref w/ path: $(/pkg-fs-ref/sub-dir)
	Interpret(
		scope map[string]*types.Value,
		expression interface{},
		opHandle types.DataHandle,
	) (*types.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		referenceInterpreter: reference.NewInterpreter(),
	}
}

type _interpreter struct {
	referenceInterpreter reference.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*types.Value,
	expression interface{},
	opHandle types.DataHandle,
) (*types.Value, error) {
	switch expression := expression.(type) {
	case string:
		// @TODO: this incorrectly treats $(inScope)$(inScope) as ref
		if strings.HasPrefix(expression, interpolater.RefStart) && strings.HasSuffix(expression, interpolater.RefEnd) {

			value, err := itp.referenceInterpreter.Interpret(
				expression,
				scope,
				opHandle,
			)
			if nil != err {
				return nil, fmt.Errorf("unable to interpret %+v to dir; error was %v", expression, err)
			}

			return value, nil

		}
	}

	return nil, fmt.Errorf("unable to interpret %v to dir", expression)

}
