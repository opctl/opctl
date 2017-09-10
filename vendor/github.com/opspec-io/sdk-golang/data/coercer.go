package data

//go:generate counterfeiter -o ./fakeCoercer.go --fake-name fakeCoercer ./ coercer

type coercer interface {
	coerceToNumber
	coerceToObject
	coerceToString
}

func newCoercer() coercer {
	return _coercer{
		coerceToNumber: newCoerceToNumber(),
		coerceToObject: newCoerceToObject(),
		coerceToString: newCoerceToString(),
	}
}

type _coercer struct {
	coerceToNumber
	coerceToObject
	coerceToString
}
