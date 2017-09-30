package outputs

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Outputs

type Outputs interface {
	interpreter
}

func New() Outputs {
	return _Outputs{
		interpreter: newInterpreter(),
	}
}

type _Outputs struct {
	interpreter
}
