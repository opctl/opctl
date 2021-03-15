package object

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
	"github.com/pkg/errors"
)

// Interpret an expression to an object value.
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
		return nil, errors.Wrap(err, fmt.Sprintf("unable to interpret %+v to object", expression))
	}

	return coerce.ToObject(&value)
}
