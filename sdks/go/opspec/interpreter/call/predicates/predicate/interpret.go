package predicate

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/eq"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/exists"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/gt"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/gte"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/lt"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/lte"
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
	case predicateSpec.LT != nil:
		return lt.Interpret(
			*predicateSpec.LT,
			scope,
		)
	case predicateSpec.LTE != nil:
		return lte.Interpret(
			*predicateSpec.LTE,
			scope,
		)
	case predicateSpec.GT != nil:
		return gt.Interpret(
			*predicateSpec.GT,
			scope,
		)
	case predicateSpec.GTE != nil:
		return gte.Interpret(
			*predicateSpec.GTE,
			scope,
		)
	case predicateSpec.NE != nil:
		return ne.Interpret(
			*predicateSpec.NE,
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
