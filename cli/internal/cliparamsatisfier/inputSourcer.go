package cliparamsatisfier

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeInputSourcer.go --fake-name FakeInputSourcer ./ InputSourcer

type InputSourcer interface {
	// Source obtains values for inputs in order of precedence.
	Source(inputName string) (*string, bool)
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
) (*string, bool) {
	for _, source := range this.sources {
		if inputValue, ok := source.ReadString(inputName); ok {
			return inputValue, true
		}
	}
	return nil, false
}
