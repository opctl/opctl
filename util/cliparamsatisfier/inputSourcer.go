package cliparamsatisfier

//go:generate counterfeiter -o ./fakeInputSourcer.go --fake-name FakeInputSourcer ./ InputSourcer

type InputSourcer interface {
	// Source obtains values for inputs in order of precedence.
	Source(inputName string) *string
}

func NewInputSourcer(
	sources ...InputSrc,
) InputSourcer {
	return inputSourcer{
		sources: sources,
	}
}

type inputSourcer struct {
	sources []InputSrc
}

func (this inputSourcer) Source(
	inputName string,
) *string {
	for _, source := range this.sources {
		if inputValue := source.Read(inputName); nil != inputValue {
			return inputValue
		}
	}
	return nil
}
