package eq

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

// Interpret an eq predicate
func Interpret(
	expressions []interface{},
	scope map[string]*model.Value,
) (bool, error) {
	var previousItemAsString string
	for i, expression := range expressions {
		// interpret items as strings since everything is coercible to string
		item, err := str.Interpret(scope, expression)
		if err != nil {
			return false, err
		}
		currentItemAsString := *item.String

		if i == 0 {
			previousItemAsString = currentItemAsString
			continue
		}

		if previousItemAsString != currentItemAsString {
			// if previous item not equal to current item, predicate is false.
			return false, nil
		}
	}
	return true, nil
}
