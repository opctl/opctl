package params

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Params

type Params interface {
	defaulter
}

func New() Params {
	return struct{ defaulter }{
		newDefaulter(),
	}
}
