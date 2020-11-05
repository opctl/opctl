package object

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
)

// Interpret interprets an expression to a object value.
// Expression must be either a type supported by coerce.ToObject
// or an object initializer
func Interpret(
	scope map[string]*model.Value,
	expression interface{},
) (*model.Value, error) {
	value, err := value.Interpret(
		expression,
		scope,
	)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret %+v to object; error was %v", expression, err)
	}

	return coerce.ToObject(&value)
}
