package str

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
)

// Interpret an expression to a string value.
func Interpret(
	scope map[string]*ipld.Node,
	expression interface{},
) (*ipld.Node, error) {
	v, err := value.Interpret(
		expression,
		scope,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to interpret %+v to string: %w", expression, err)
	}

	return coerce.ToString(&v)
}
