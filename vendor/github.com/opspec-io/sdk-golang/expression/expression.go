// Package expression implements use cases surrounding expressions such as evaluation
package expression

import (
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"strings"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Expression

type Expression interface {
	evalArrayer
	evalDirer
	evalFiler
	evalNumberer
	evalObjecter
	evalStringer
}

func New() Expression {
	return _Expression{
		evalArrayer:  newEvalArrayer(),
		evalDirer:    newEvalDirer(),
		evalFiler:    newEvalFiler(),
		evalNumberer: newEvalNumberer(),
		evalObjecter: newEvalObjecter(),
		evalStringer: newEvalStringer(),
	}
}

type _Expression struct {
	evalArrayer
	evalDirer
	evalFiler
	evalNumberer
	evalObjecter
	evalStringer
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
