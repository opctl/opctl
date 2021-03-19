package notexists

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
)

//Interpret a notExists expression
func Interpret(
	expression string,
	scope map[string]*model.Value,
) (bool, error) {
	// @TODO: make more exact. reference.Interpret can err for more reasons than simply null pointer exceptions.
	_, err := reference.Interpret(
		expression,
		scope,
		nil,
	)

	return err != nil, nil
}
