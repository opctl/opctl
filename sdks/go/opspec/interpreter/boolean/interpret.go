package boolean

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
	"github.com/pkg/errors"
)

// Interpret an expression to a boolean value.
// Expression must be a type supported by coerce.ToBoolean
func Interpret(
	scope map[string]*model.Value,
	expression interface{},
) (*model.Value, error) {
	v, err := value.Interpret(
		expression,
		scope,
	)
	if nil != err {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to interpret %+v to boolean", expression))
	}

	return coerce.ToBoolean(&v)
}
