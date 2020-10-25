package predicate

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/eq"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/exists"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/ne"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/notexists"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	Interpret(
		predicateSpec *model.PredicateSpec,
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
	predicateSpec *model.PredicateSpec,
	scope map[string]*model.Value,
) (bool, error) {
	switch {
	case nil != predicateSpec.Eq:
		return itp.eqInterpreter.Interpret(
			*predicateSpec.Eq,
			scope,
		)
	case nil != predicateSpec.Exists:
		return itp.existsInterpreter.Interpret(
			*predicateSpec.Exists,
			scope,
		)
	case nil != predicateSpec.Ne:
		return itp.neInterpreter.Interpret(
			*predicateSpec.Ne,
			scope,
		)
	case nil != predicateSpec.NotExists:
		return itp.notExistsInterpreter.Interpret(
			*predicateSpec.NotExists,
			scope,
		)
	default:
		return false, fmt.Errorf("unable to interpret predicate; predicate was unexpected type %+v", predicateSpec)
	}
}
