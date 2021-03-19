package ne

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

// Interpret a ne expression
func Interpret(
	expressions []interface{},
	scope map[string]*model.Value,
) (bool, error) {
	var itemsAsStrings []string
	for _, expression := range expressions {
		// interpret items as strings since everything is coercible to string
		item, err := str.Interpret(scope, expression)
		if err != nil {
			return false, err
		}
		currentItemAsString := *item.String

		for _, prevItemAsString := range itemsAsStrings {
			// if any previously seen item equals current item predicate is false.
			if prevItemAsString == currentItemAsString {
				return false, nil
			}
		}

		itemsAsStrings = append(itemsAsStrings, currentItemAsString)
	}
	return true, nil
}
