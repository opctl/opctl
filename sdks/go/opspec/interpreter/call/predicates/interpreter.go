package predicates

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	Interpret(
		predicateSpecs []*model.PredicateSpec,
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
	predicateSpecs []*model.PredicateSpec,
	scope map[string]*model.Value,
) (bool, error) {
	for _, predicateSpec := range predicateSpecs {
		dcgPredicate, err := itp.predicateInterpreter.Interpret(
			predicateSpec,
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
