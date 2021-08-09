package array

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
)

// Interpret an expression to an array value.
// Expression must be either a type supported by coerce.ToArray
// or an array initializer
func Interpret(
	scope map[string]*model.Value,
	expression interface{},
) (*model.Value, error) {
	v, err := value.Interpret(
		expression,
		scope,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to interpret %+v to array: %w", expression, err)
	}

	return coerce.ToArray(&v)
}
