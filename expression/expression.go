// Package expression implements use cases surrounding expressions such as evaluation
package expression

import (
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"strings"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Expression

type Expression interface {
	dirEvaluator
	fileEvaluator
	numberEvaluator
	objectEvaluator
	stringEvaluator
}

func New() Expression {
	return _Expression{
		dirEvaluator:    newDirEvaluator(),
		fileEvaluator:   newFileEvaluator(),
		numberEvaluator: newNumberEvaluator(),
		objectEvaluator: newObjectEvaluator(),
		stringEvaluator: newStringEvaluator(),
	}
}

type _Expression struct {
	dirEvaluator
	fileEvaluator
	numberEvaluator
	objectEvaluator
	stringEvaluator
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
