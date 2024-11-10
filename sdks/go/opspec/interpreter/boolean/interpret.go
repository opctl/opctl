package boolean

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
)

// Interpret an expression to a boolean value.
// Expression must be a type supported by coerce.ToBoolean
func Interpret(
	scope map[string]*ipld.Node,
	expression interface{},
) (*ipld.Node, error) {
	v, err := value.Interpret(
		expression,
		scope,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to interpret %+v to boolean: %w", expression, err)
	}

	return coerce.ToBoolean(&v)
}
