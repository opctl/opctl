package str

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
	"github.com/pkg/errors"
)

// Interpret an expression to a string value.
func Interpret(
	scope map[string]*model.Value,
	expression interface{},
) (*model.Value, error) {
	v, err := value.Interpret(
		expression,
		scope,
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to interpret %+v to string", expression))
	}

	return coerce.ToString(&v)
}
