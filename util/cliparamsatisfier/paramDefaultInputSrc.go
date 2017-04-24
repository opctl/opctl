package cliparamsatisfier

import (
	"github.com/opspec-io/sdk-golang/model"
	"strconv"
)

func NewParamDefaultInputSrc(
	inputs map[string]*model.Param,
) InputSrc {
	return paramDefaultInputSrc{
		inputs:      inputs,
		readHistory: map[string]struct{}{},
	}
}

// paramDefaultInputSrc implements InputSrc interface by sourcing inputs from input defaults
type paramDefaultInputSrc struct {
	inputs      map[string]*model.Param
	readHistory map[string]struct{} // tracks reads
}

func (this paramDefaultInputSrc) Read(
	inputName string,
) *string {
	if _, ok := this.readHistory[inputName]; ok {
		// enforce read at most once.
		return nil
	}

	if inputValue, ok := this.inputs[inputName]; ok {
		// track read history
		this.readHistory[inputName] = struct{}{}

		switch {
		case nil != inputValue.Dir:
			return inputValue.Dir.Default
		case nil != inputValue.File:
			return inputValue.File.Default
		case nil != inputValue.Number && nil != inputValue.Number.Default:
			floatString := strconv.FormatFloat(*inputValue.Number.Default, 'E', -1, 64)
			return &floatString
		case nil != inputValue.String:
			return inputValue.String.Default
		}
	}
	return nil
}
