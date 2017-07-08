package cliparamsatisfier

import (
	"strings"
)

func (isf _InputSrcFactory) NewSliceInputSrc(
	args []string,
	sep string,
) InputSrc {
	argMap := map[string]string{}
	for _, arg := range args {
		// get parts
		parts := strings.SplitN(arg, sep, 2)
		inputName := parts[0]
		inputValue := parts[1]

		argMap[inputName] = inputValue
	}
	return sliceInputSrc{argMap}
}

// sliceInputSrc implements InputSrc interface by sourcing inputs from a slice
type sliceInputSrc struct {
	argMap map[string]string
}

func (this sliceInputSrc) ReadString(
	inputName string,
) (*string, bool) {
	if inputValue, ok := this.argMap[inputName]; ok {
		// enforce read at most once.
		delete(this.argMap, inputName)

		return &inputValue, true
	}
	return nil, false
}
