package interpolater

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Interpolater

import (
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater/dereferencer"
)

const (
	escaper   = '\\'
	operator  = '$'
	refOpener = '('
	refCloser = ')'
	RefStart  = string(operator) + string(refOpener)
	RefEnd    = string(refCloser)
)

type Interpolater interface {
	// Interpolate interpolates the provided expression
	Interpolate(
		expression string,
		scope map[string]*model.Value,
		opHandle model.DataHandle,
	) (string, error)
}

func New() Interpolater {
	return _Interpolater{
		deReferencer: dereferencer.New(),
	}
}

type _Interpolater struct {
	deReferencer dereferencer.DeReferencer
}

// similar: https://github.com/kubernetes/kubernetes/blob/5066a67caaf8638c7473d4bd228037d0c270c546/third_party/forked/golang/expansion/expand.go#L1
func (itp _Interpolater) Interpolate(
	expression string,
	scope map[string]*model.Value,
	opHandle model.DataHandle,
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

			result, consumed, err := itp.tryDeRef(expression[i+1:], scope, opHandle)
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
// returns the de referenced value (if any), number of bytes consumed, and any err
func (itp _Interpolater) tryDeRef(
	possibleRef string,
	scope map[string]*model.Value,
	opHandle model.DataHandle,
) (string, int, error) {
	refBuffer := []byte{}
	i := 0

	for i < len(possibleRef) {
		switch possibleRef[i] {
		case refCloser:
			if len(refBuffer) > 0 && refOpener == refBuffer[0] {
				result, ok, err := itp.deReferencer.DeReference(string(refBuffer[1:]), scope, opHandle)
				if nil != err {
					return "", 0, err
				}
				if ok {
					return result, i + 1, err
				}
			}
			refBuffer = append(refBuffer, possibleRef[i])
		case operator:
			result, consumed, err := itp.tryDeRef(possibleRef[i+1:], scope, opHandle)
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
