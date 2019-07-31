package predicates

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

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
