package lt

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/number"
)

// Interpret a lt predicate
func Interpret(
	expressions []interface{},
	scope map[string]*model.Value,
) (bool, error) {
	var previousItemAsNumber float64
	for i, expression := range expressions {
		item, err := number.Interpret(scope, expression)
		if err != nil {
			return false, err
		}
		currentItemAsNumber := *item.Number

		if i == 0 {
			previousItemAsNumber = currentItemAsNumber
			continue
		}

		if !(previousItemAsNumber < currentItemAsNumber) {
			// if previous item not ltss than the current item, predicate is false.
			return false, nil
		}
	}
	return true, nil
}
