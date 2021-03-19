package interpolater

import (
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
)

const (
	escaper   = '\\'
	operator  = '$'
	refOpener = '('
	refCloser = ')'
	RefStart  = string(operator) + string(refOpener)
	RefEnd    = string(refCloser)
)

// Interpolate interpolates the provided expression
// similar: https://github.com/kubernetes/kubernetes/blob/5066a67caaf8638c7473d4bd228037d0c270c546/third_party/forked/golang/expansion/expand.go#L1
func Interpolate(
	expression string,
	scope map[string]*model.Value,
) (string, error) {
	refBuffer := []byte{}
	i := 0
	escapesCount := 0

	for i < len(expression) {
		switch expression[i] {
		case escaper:
			escapesCount++
		case operator:
			isEscaped := 0 != escapesCount%2
			for escapesCount > 0 {
				if 0 == escapesCount%2 {
					refBuffer = append(refBuffer, escaper)
				}
				escapesCount--
			}
			if isEscaped {
				refBuffer = append(refBuffer, expression[i])
				break
			}

			result, consumed, err := tryDeRef(expression[i+1:], scope)
			if nil != err {
				return "", err
			}
			refBuffer = append(refBuffer, result...)
			i += consumed
		default:
			for escapesCount > 0 {
				refBuffer = append(refBuffer, escaper)
				escapesCount--
			}

			refBuffer = append(refBuffer, expression[i])
		}

		// always increment loop counter
		i++
	}

	for escapesCount > 0 {
		refBuffer = append(refBuffer, escaper)
		escapesCount--
	}

	return string(refBuffer), nil
}

// tryDeRef tries to de reference from possibleRef.
// It's assumed possibleRef doesn't contain an initial operator.
//
// returns the interpreted value (if any), number of bytes consumed, and any err
func tryDeRef(
	possibleRef string,
	scope map[string]*model.Value,
) (string, int, error) {
	refBuffer := []byte{}
	i := 0

	for i < len(possibleRef) {
		switch possibleRef[i] {
		case refCloser:
			if len(refBuffer) > 0 && refOpener == refBuffer[0] {
				value, err := reference.Interpret(opspec.NameToRef(string(refBuffer[1:])), scope, nil)
				if nil != err {
					return "", 0, err
				}

				valueAsString, err := coerce.ToString(value)
				if nil != err {
					return "", 0, err
				}

				return *valueAsString.String, i + 1, err
			}
			refBuffer = append(refBuffer, possibleRef[i])
		case operator:
			result, consumed, err := tryDeRef(possibleRef[i+1:], scope)
			if nil != err {
				return "", 0, err
			}
			refBuffer = append(refBuffer, result...)
			i += consumed
		default:
			refBuffer = append(refBuffer, possibleRef[i])
		}

		// always increment loop counter
		i++
	}

	return "$" + string(refBuffer), len(possibleRef), nil
}
