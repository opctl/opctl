// Package expression implements use cases surrounding expressions such as evaluation
package expression

import (
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"strings"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Expression

type Expression interface {
	evalToDir
	evalToFile
	evalToNumber
	evalToObject
	evalToString
}

func New() Expression {
	return _Expression{
		evalToDir:    newEvalToDir(),
		evalToFile:   newEvalToFile(),
		evalToNumber: newEvalToNumber(),
		evalToObject: newEvalToObject(),
		evalToString: newEvalToString(),
	}
}

type _Expression struct {
	evalToDir
	evalToFile
	evalToNumber
	evalToObject
	evalToString
}

func tryResolveExplicitRef(
	expression string,
	scope map[string]*model.Value,
) (*model.Value, bool) {
	possibleRef := strings.TrimPrefix(expression, string(interpolater.Operator+interpolater.RefOpener))
	possibleRef = strings.TrimSuffix(possibleRef, string(interpolater.RefCloser))
	dcgValue, ok := scope[possibleRef]
	return dcgValue, ok
}
