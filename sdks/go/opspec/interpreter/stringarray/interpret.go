package stringarray

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/array"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

// Interpret a string array
func Interpret(
	scope map[string]*model.Value,
	expression interface{},
) ([]string, error) {
	if expression == nil {
		return []string{}, nil
	}

	a, err := array.Interpret(
		scope,
		expression,
	)
	if err != nil {
		return nil, err
	}

	sa := []string{}

	for _, cmdEntryExpression := range *a.Array {
		// interpret each entry as string
		cmdEntry, err := str.Interpret(scope, cmdEntryExpression)
		if err != nil {
			return nil, err
		}
		sa = append(sa, *cmdEntry.String)
	}

	return sa, nil
}
