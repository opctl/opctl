package predicates

import (
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate"
	"github.com/opctl/opctl/sdks/go/types"
)

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

type Interpreter interface {
	Interpret(
		opHandle types.DataHandle,
		scgPredicates []*types.SCGPredicate,
		scope map[string]*types.Value,
	) (bool, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return &_interpreter{
		predicateInterpreter: predicate.NewInterpreter(),
	}
}

type _interpreter struct {
	predicateInterpreter predicate.Interpreter
}

func (itp _interpreter) Interpret(
	opHandle types.DataHandle,
	scgPredicates []*types.SCGPredicate,
	scope map[string]*types.Value,
) (bool, error) {
	for _, scgPredicate := range scgPredicates {
		dcgPredicate, err := itp.predicateInterpreter.Interpret(
			opHandle,
			scgPredicate,
			scope,
		)
		if nil != err {
			return false, err
		}

		if !dcgPredicate {
			return false, nil
		}
	}
	return true, nil
}
