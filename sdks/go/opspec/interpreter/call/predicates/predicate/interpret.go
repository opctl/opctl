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
	case predicateSpec.Eq != nil:
		return eq.Interpret(
			*predicateSpec.Eq,
			scope,
		)
	case predicateSpec.Exists != nil:
		return exists.Interpret(
			*predicateSpec.Exists,
			scope,
		)
	case predicateSpec.Ne != nil:
		return ne.Interpret(
			*predicateSpec.Ne,
			scope,
		)
	case predicateSpec.NotExists != nil:
		return notexists.Interpret(
			*predicateSpec.NotExists,
			scope,
		)
	default:
		return false, fmt.Errorf("unable to interpret predicate: predicate was unexpected type %+v", predicateSpec)
	}
}
