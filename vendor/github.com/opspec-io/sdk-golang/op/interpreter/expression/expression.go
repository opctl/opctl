// Package expression implements use cases surrounding expressions such as evaluation
package expression

import (
	"strings"

	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression/interpolater"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Expression

type Expression interface {
	evalArrayer
	evalBooleaner
	evalDirer
	evalFiler
	evalNumberer
	evalObjecter
	evalStringer
}

func New() Expression {
	return struct {
		evalArrayer
		evalBooleaner
		evalDirer
		evalFiler
		evalNumberer
		evalObjecter
		evalStringer
	}{
		evalArrayer:   newEvalArrayer(),
		evalBooleaner: newEvalBooleaner(),
		evalDirer:     newEvalDirer(),
		evalFiler:     newEvalFiler(),
		evalNumberer:  newEvalNumberer(),
		evalObjecter:  newEvalObjecter(),
		evalStringer:  newEvalStringer(),
	}
}

func tryResolveExplicitRef(
	expression string,
	scope map[string]*model.Value,
) (*model.Value, bool) {
	if strings.HasPrefix(expression, interpolater.RefStart) && strings.HasSuffix(expression, interpolater.RefEnd) {
		dcgValue, ok := scope[expression[2:len(expression)-1]]
		return dcgValue, ok
	}

	return nil, false
}
