package interpolater

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Interpolater
//go:generate counterfeiter -o ./fakeDeReferencer.go --fake-name FakeDeReferencer ./ DeReferencer

const (
	operator  = '$'
	refOpener = '('
	refCloser = ')'
)

// ValueSourcer de references references
type DeReferencer interface {
	DeReference(
		ref string,
	) (string, error)
}

type Interpolater interface {
	Interpolate(
		expression string,
		deReferencer DeReferencer,
	) (string, error)
}

func New() Interpolater {
	return _Interpolater{}
}

type _Interpolater struct{}

func (itp _Interpolater) Interpolate(
	expression string,
	deReferencer DeReferencer,
) (string, error) {
	refBuffer := []byte{}
	i := 0
	for i < len(expression) {
		switch {
		case operator == expression[i]:
			value, consumed, ok, err := itp.tryDeRef(expression[i+1:], deReferencer)
			if nil != err {
				return "", err
			}
			if ok {
				refBuffer = append(refBuffer, value...)
			} else {
				refBuffer = append(refBuffer, expression[i:i+consumed+1]...)
			}
			i += consumed
		default:
			refBuffer = append(refBuffer, expression[i])
		}

		// always increment loop counter
		i++
	}

	return string(refBuffer), nil
}

// tryDeRef tries to de reference from possibleRef.
// It's assumed possibleRef doesn't contain the initial operator.
//
// returns the de referenced value (if any), number of bytes consumed, whether de referencing occurred, and any err
func (itp _Interpolater) tryDeRef(
	possibleRef string,
	deReferencer DeReferencer,
) (string, int, bool, error) {
	refBuffer := []byte{}
	i := 0
	for i < len(possibleRef) {
		switch {
		case i == 0:
			if refOpener != possibleRef[i] {
				return "", 1, false, nil
			}
		case refCloser == possibleRef[i]:
			value, err := deReferencer.DeReference(string(refBuffer))
			return value, i + 1, true, err
		case operator == possibleRef[i]:
			value, consumed, ok, err := itp.tryDeRef(possibleRef[i+1:], deReferencer)
			if nil != err {
				return "", 0, false, err
			}
			if ok {
				refBuffer = append(refBuffer, value...)
			} else {
				refBuffer = append(refBuffer, possibleRef[i:i+consumed]...)
			}
			i += consumed
		default:
			refBuffer = append(refBuffer, possibleRef[i])
		}

		// always increment loop counter
		i++
	}

	return "", len(possibleRef), false, nil
}
