package number

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name fakeInterpreter ./ Interpreter

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"strconv"
)

type Interpreter interface {
	// interprets an expression to a string
	Interpret(
		scope map[string]*model.Value,
		expression string,
	) (float64, error)
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
) (float64, error) {
	var resultBuffer, possibleRefBuffer bytes.Buffer

	// note: WriteByte/Number errs ignored as per their docs; they're always nil
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
				return 0, fmt.Errorf("Unable to interpret number; %v not in scope", ref)
			}

			numberValue, err := itp.coercer.Coerce(value)
			if nil != err {
				return 0, fmt.Errorf("Unable to interpret number; error was: %v", err.Error())
			}

			resultBuffer.WriteString(strconv.FormatFloat(numberValue, 'f', -1, 64))
		case possibleRefBuffer.Len() > 0:
			possibleRefBuffer.WriteByte(expression[i])
		default:
			resultBuffer.WriteByte(expression[i])
		}
	}

	float64Value, err := strconv.ParseFloat(resultBuffer.String()+possibleRefBuffer.String(), 64)
	if nil != err {
		return 0, fmt.Errorf("Unable to interpret number; error was %v", err.Error())
	}
	return float64Value, nil
}
