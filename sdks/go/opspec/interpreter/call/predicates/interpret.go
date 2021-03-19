package predicates

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate"
)

// Interpret predicates
func Interpret(
	predicateSpecs []*model.PredicateSpec,
	scope map[string]*model.Value,
) (bool, error) {
	for _, predicateSpec := range predicateSpecs {
		predicate, err := predicate.Interpret(
			predicateSpec,
			scope,
		)
		if err != nil {
			return false, err
		}

		if !predicate {
			return false, nil
		}
	}
	return true, nil
}
