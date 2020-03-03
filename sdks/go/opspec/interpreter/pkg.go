// Package interpreter defines an interpreter for ops.
//
// Implementation
//
// The interpreter is a Tree-Walking interpreter.
// Values are represented by the Value struct of the github.com/opctl/opctl/sdks/go/model package.
package interpreter

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater"
	"strings"
)

// @TODO: find this a better home
func TryResolveExplicitRef(
	expression string,
	scope map[string]*model.Value,
) (*model.Value, bool) {
	if strings.HasPrefix(expression, interpolater.RefStart) && strings.HasSuffix(expression, interpolater.RefEnd) {
		dcgValue, ok := scope[expression[2:len(expression)-1]]
		return dcgValue, ok
	}

	return nil, false
}
