package number

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
)

// Interpret an expression to a number value.
// Expression must be either a type supported by coerce.ToNumber
// or an number initializer
func Interpret(
	scope map[string]*model.Value,
	expression interface{},
) (*model.Value, error) {
	v, err := value.Interpret(
		expression,
		scope,
	)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret %+v to number; error was %v", expression, err)
	}

	return coerce.ToNumber(&v)
}
