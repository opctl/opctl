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
	// DeReference returns the de referenced value (if any), whether de referencing occurred, and any err
	DeReference(
		ref string,
	) (string, bool, error)
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
			result, consumed, err := itp.tryDeRef(expression[i+1:], deReferencer)
			if nil != err {
				return "", err
			}
			refBuffer = append(refBuffer, result...)
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
) (string, int, error) {
	refBuffer := []byte{}
	i := 0
	for i < len(possibleRef) {
		switch {
		case refCloser == possibleRef[i]:
			if len(refBuffer) > 0 && refOpener == refBuffer[0] {
				result, ok, err := deReferencer.DeReference(string(refBuffer[1:]))
				if nil != err {
					return "", 0, err
				}
				if ok {
					return result, i + 1, err
				}
			}
			refBuffer = append(refBuffer, possibleRef[i])
		case operator == possibleRef[i]:
			result, consumed, err := itp.tryDeRef(possibleRef[i+1:], deReferencer)
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
