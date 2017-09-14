package interpolater

import (
	"github.com/opspec-io/sdk-golang/model"
	"strings"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Interpolater

const (
	operator  = '$'
	refOpener = '('
	refCloser = ')'
)

type Interpolater interface {
	Interpolate(
		expression string,
		scope map[string]*model.Value,
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func New() Interpolater {
	return _Interpolater{
		deReferencer: newDeReferencer(),
	}
}

type _Interpolater struct {
	deReferencer
}

func (itp _Interpolater) Interpolate(
	expression string,
	scope map[string]*model.Value,
	pkgHandle model.PkgHandle,
) (*model.Value, error) {

	if dcgValue, ok := scope[strings.TrimSuffix(strings.TrimPrefix(expression, "$("), ")")]; ok {
		// explicit arg
		return dcgValue, nil
	}

	refBuffer := []byte{}
	i := 0
	for i < len(expression) {
		switch {
		case operator == expression[i]:
			result, consumed, err := itp.tryDeRef(expression[i+1:], scope, pkgHandle)
			if nil != err {
				return nil, err
			}
			refBuffer = append(refBuffer, result...)
			i += consumed
		default:
			refBuffer = append(refBuffer, expression[i])
		}

		// always increment loop counter
		i++
	}

	valueAsString := string(refBuffer)
	return &model.Value{String: &valueAsString}, nil
}

// tryDeRef tries to de reference from possibleRef.
// It's assumed possibleRef doesn't contain the initial operator.
//
// returns the de referenced value (if any), number of bytes consumed, and any err
func (itp _Interpolater) tryDeRef(
	possibleRef string,
	scope map[string]*model.Value,
	pkgHandle model.PkgHandle,
) (string, int, error) {
	refBuffer := []byte{}
	i := 0
	for i < len(possibleRef) {
		switch {
		case refCloser == possibleRef[i]:
			if len(refBuffer) > 0 && refOpener == refBuffer[0] {
				result, ok, err := itp.deReferencer.DeReference(string(refBuffer[1:]), scope, pkgHandle)
				if nil != err {
					return "", 0, err
				}
				if ok {
					return result, i + 1, err
				}
			}
			refBuffer = append(refBuffer, possibleRef[i])
		case operator == possibleRef[i]:
			result, consumed, err := itp.tryDeRef(possibleRef[i+1:], scope, pkgHandle)
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
