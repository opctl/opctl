package data

//go:generate counterfeiter -o ./fakeCoercer.go --fake-name fakeCoercer ./ coercer

type coercer interface {
	coerceToArray
	coerceToBoolean
	coerceToFile
	coerceToNumber
	coerceToObject
	coerceToString
}

func newCoercer() coercer {
	return struct {
		coerceToArray
		coerceToBoolean
		coerceToFile
		coerceToNumber
		coerceToObject
		coerceToString
	}{
		newCoerceToArray(),
		newCoerceToBoolean(),
		newCoerceToFile(),
		newCoerceToNumber(),
		newCoerceToObject(),
		newCoerceToString(),
	}
}
