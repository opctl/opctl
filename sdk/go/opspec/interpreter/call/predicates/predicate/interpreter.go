package predicate

import (
	"fmt"

	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/call/predicates/predicate/eq"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/call/predicates/predicate/exists"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/call/predicates/predicate/ne"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/call/predicates/predicate/notexists"
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
		eqInterpreter:        eq.NewInterpreter(),
		existsInterpreter:    exists.NewInterpreter(),
		neInterpreter:        ne.NewInterpreter(),
		notExistsInterpreter: notexists.NewInterpreter(),
	}
}

type _interpreter struct {
	eqInterpreter        eq.Interpreter
	existsInterpreter    exists.Interpreter
	neInterpreter        ne.Interpreter
	notExistsInterpreter notexists.Interpreter
}

func (itp _interpreter) Interpret(
	opHandle model.DataHandle,
	scgPredicate *model.SCGPredicate,
	scope map[string]*model.Value,
) (bool, error) {
	switch {
	case nil != scgPredicate.Eq:
		return itp.eqInterpreter.Interpret(
			*scgPredicate.Eq,
			opHandle,
			scope,
		)
	case nil != scgPredicate.Exists:
		return itp.existsInterpreter.Interpret(
			*scgPredicate.Exists,
			opHandle,
			scope,
		)
	case nil != scgPredicate.Ne:
		return itp.neInterpreter.Interpret(
			*scgPredicate.Ne,
			opHandle,
			scope,
		)
	case nil != scgPredicate.NotExists:
		return itp.notExistsInterpreter.Interpret(
			*scgPredicate.NotExists,
			opHandle,
			scope,
		)
	default:
		return false, fmt.Errorf("unable to interpret predicate; predicate was unexpected type %+v", scgPredicate)
	}
}
