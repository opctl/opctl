package data

//go:generate counterfeiter -o ./fakeCoercer.go --fake-name fakeCoercer ./ coercer

type coercer interface {
	coerceToArray
	coerceToFile
	coerceToNumber
	coerceToObject
	coerceToString
}

func newCoercer() coercer {
	return _coercer{
		coerceToArray:  newCoerceToArray(),
		coerceToFile:   newCoerceToFile(),
		coerceToNumber: newCoerceToNumber(),
		coerceToObject: newCoerceToObject(),
		coerceToString: newCoerceToString(),
	}
}

type _coercer struct {
	coerceToArray
	coerceToFile
	coerceToNumber
	coerceToObject
	coerceToString
}
