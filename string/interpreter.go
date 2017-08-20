package string

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name fakeInterpreter ./ Interpreter

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
)

type Interpreter interface {
	// interprets an expression to a string
	Interpret(
		scope map[string]*model.Value,
		expression string,
	) (string, error)
}

func newInterpreter() Interpreter {
	return _interpreter{
		coercer: newCoercer(),
	}
}

type _interpreter struct {
	coercer coercer
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	expression string,
) (string, error) {
	var resultBuffer, possibleRefBuffer bytes.Buffer

	// note: WriteByte/String errs ignored as per their docs; they're always nil
	for i := 0; i < len(expression); i++ {
		switch {
		case '$' == expression[i]:
			possibleRefBuffer.WriteByte('$')
		case possibleRefBuffer.Len() == 1 && '(' == expression[i]:
			possibleRefBuffer.WriteByte('(')
		case possibleRefBuffer.Len() > 0 && ')' == expression[i]:
			// we've got a ref
			ref := possibleRefBuffer.String()[2:]
			possibleRefBuffer.Reset()

			value, ok := scope[ref]
			if !ok {
				return "", fmt.Errorf("Unable to interpret string; %v not in scope", ref)
			}

			stringValue, err := itp.coercer.Coerce(value)
			if nil != err {
				return "", fmt.Errorf("Unable to interpret string; error was: %v", err.Error())
			}

			resultBuffer.WriteString(stringValue)
		case possibleRefBuffer.Len() > 0:
			possibleRefBuffer.WriteByte(expression[i])
		default:
			resultBuffer.WriteByte(expression[i])
		}
	}

	return resultBuffer.String() + possibleRefBuffer.String(), nil
}
