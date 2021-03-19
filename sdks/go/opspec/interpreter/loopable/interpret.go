package loopable

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/array"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/object"
)

//Interpret an expression to a loopable value.
func Interpret(
	expression interface{},
	scope map[string]*model.Value,
) (*model.Value, error) {
	// try interpreting as array
	if dcgForEach, err := array.Interpret(
		scope,
		expression,
	); err == nil {
		return dcgForEach, err
	}

	// fallback to interpreting as object
	return object.Interpret(
		scope,
		expression,
	)
}
