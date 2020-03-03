package envvar

import (
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc"
)

func New() inputsrc.InputSrc {
	return envVarInputSrc{
		os:          ios.New(),
		readHistory: map[string]struct{}{},
	}
}

// envVarInputSrc implements InputSrc interface by sourcing inputs from env vars
type envVarInputSrc struct {
	os          ios.IOS
	readHistory map[string]struct{} // tracks reads
}

func (this envVarInputSrc) ReadString(
	inputName string,
) (*string, bool) {
	if _, ok := this.readHistory[inputName]; ok {
		// enforce read at most once.
		return nil, false
	}

	if inputValue := this.os.Getenv(inputName); "" != inputValue {
		// track read history
		this.readHistory[inputName] = struct{}{}

		return &inputValue, true
	}
	return nil, false
}
