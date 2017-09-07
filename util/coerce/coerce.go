package coerce

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Coerce

type Coerce interface {
	toNumber
	toString
}

func New() Coerce {
	return _Coerce{
		toNumber: newToNumber(),
		toObject: newToObject(),
		toString: newToString(),
	}
}

type _Coerce struct {
	toNumber
	toObject
	toString
}
