package predicate

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/eq"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/exists"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/ne"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/notexists"
)

// Interpret a predicate
func Interpret(
	predicateSpec *model.PredicateSpec,
	scope map[string]*model.Value,
) (bool, error) {
	switch {
	case nil != predicateSpec.Eq:
		return eq.Interpret(
			*predicateSpec.Eq,
			scope,
		)
	case nil != predicateSpec.Exists:
		return exists.Interpret(
			*predicateSpec.Exists,
			scope,
		)
	case nil != predicateSpec.Ne:
		return ne.Interpret(
			*predicateSpec.Ne,
			scope,
		)
	case nil != predicateSpec.NotExists:
		return notexists.Interpret(
			*predicateSpec.NotExists,
			scope,
		)
	default:
		return false, fmt.Errorf("unable to interpret predicate: predicate was unexpected type %+v", predicateSpec)
	}
}
