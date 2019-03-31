package predicate

import (
	"fmt"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/predicates/predicate/eq"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/predicates/predicate/ne"
)

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

type Interpreter interface {
	Interpret(
		opHandle model.DataHandle,
		scgPredicate *model.SCGPredicate,
		scope map[string]*model.Value,
	) (bool, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return &_interpreter{
		eqInterpreter: eq.NewInterpreter(),
		neInterpreter: ne.NewInterpreter(),
	}
}

type _interpreter struct {
	eqInterpreter eq.Interpreter
	neInterpreter ne.Interpreter
}

func (itp _interpreter) Interpret(
	opHandle model.DataHandle,
	scgPredicate *model.SCGPredicate,
	scope map[string]*model.Value,
) (bool, error) {
	switch {
	case nil != scgPredicate.Eq:
		return itp.eqInterpreter.Interpret(
			scgPredicate.Eq,
			opHandle,
			scope,
		)
	case nil != scgPredicate.Ne:
		return itp.neInterpreter.Interpret(
			scgPredicate.Ne,
			opHandle,
			scope,
		)
	default:
		return false, fmt.Errorf("unable to interpret predicate; predicate was unexpected type %+v", scgPredicate)
	}
}
