package object

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
)

// Interpret an expression to an object value.
// Expression must be either a type supported by coerce.ToObject
// or an object initializer
func Interpret(
	scope map[string]*ipld.Node,
	expression interface{},
) (*ipld.Node, error) {
	value, err := value.Interpret(
		expression,
		scope,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to interpret %+v to object: %w", expression, err)
	}

	return coerce.ToObject(&value)
}
