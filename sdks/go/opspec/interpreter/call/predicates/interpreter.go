package predicates

import (
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/predicates/predicate"
)

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

type Interpreter interface {
	Interpret(
		opHandle model.DataHandle,
		scgPredicates []*model.SCGPredicate,
		scope map[string]*model.Value,
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
	opHandle model.DataHandle,
	scgPredicates []*model.SCGPredicate,
	scope map[string]*model.Value,
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
