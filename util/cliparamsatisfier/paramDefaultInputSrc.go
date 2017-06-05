package cliparamsatisfier

import (
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"strconv"
	"strings"
)

func NewParamDefaultInputSrc(
	inputs map[string]*model.Param,
	pkgPath string,
) InputSrc {
	return paramDefaultInputSrc{
		inputs:      inputs,
		readHistory: map[string]struct{}{},
		pkgPath:     pkgPath,
	}
}

// paramDefaultInputSrc implements InputSrc interface by sourcing inputs from input defaults
type paramDefaultInputSrc struct {
	inputs      map[string]*model.Param
	readHistory map[string]struct{} // tracks reads
	pkgPath     string
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
		case nil != inputValue.Dir && nil != inputValue.Dir.Default:
			if strings.HasPrefix(*inputValue.Dir.Default, "/") {
				// defaulted to pkg dir
				value := filepath.Join(this.pkgPath, *inputValue.Dir.Default)
				return &value
			}
			return inputValue.Dir.Default
		case nil != inputValue.File && nil != inputValue.File.Default:
			if strings.HasPrefix(*inputValue.File.Default, "/") {
				// defaulted to pkg file
				value := filepath.Join(this.pkgPath, *inputValue.File.Default)
				return &value
			}
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
