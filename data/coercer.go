package data

//go:generate counterfeiter -o ./fakeCoercer.go --fake-name fakeCoercer ./ coercer

type coercer interface {
	coerceToFile
	coerceToNumber
	coerceToObject
	coerceToString
}

func newCoercer() coercer {
	return _coercer{
		coerceToFile:   newCoerceToFile(),
		coerceToNumber: newCoerceToNumber(),
		coerceToObject: newCoerceToObject(),
		coerceToString: newCoerceToString(),
	}
}

type _coercer struct {
	coerceToFile
	coerceToNumber
	coerceToObject
	coerceToString
}
