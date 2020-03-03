package cliparamsatisfier

import (
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc"
)

//counterfeiter:generate -o fakes/inputSourcer.go . InputSourcer
type InputSourcer interface {
	inputSourcer
}

// inputSourcer is an internal version of Parser so fakes don't cause cyclic deps
//counterfeiter:generate -o internal/fakes/InputSourcer.go . inputSourcer
type inputSourcer interface {
	// Source obtains values for inputs in order of precedence.
	Source(inputName string) (*string, bool)
}

func NewInputSourcer(
	sources ...inputsrc.InputSrc,
) InputSourcer {
	return _inputSourcer{
		sources: sources,
	}
}

type _inputSourcer struct {
	sources []inputsrc.InputSrc
}

func (this _inputSourcer) Source(
	inputName string,
) (*string, bool) {
	for _, source := range this.sources {
		if inputValue, ok := source.ReadString(inputName); ok {
			return inputValue, true
		}
	}
	return nil, false
}
