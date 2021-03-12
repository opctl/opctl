package envvar

import (
	"os"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc"
)

func New() inputsrc.InputSrc {
	return envVarInputSrc{
		readHistory: map[string]struct{}{},
	}
}

// envVarInputSrc implements InputSrc interface by sourcing inputs from env vars
type envVarInputSrc struct {
	readHistory map[string]struct{} // tracks reads
}

func (this envVarInputSrc) ReadString(
	inputName string,
) (*string, bool) {
	if _, ok := this.readHistory[inputName]; ok {
		// enforce read at most once.
		return nil, false
	}

	if inputValue := os.Getenv(inputName); "" != inputValue {
		// track read history
		this.readHistory[inputName] = struct{}{}

		return &inputValue, true
	}
	return nil, false
}
