package cliparamsatisfier

import (
	"github.com/virtual-go/fs/osfs"
	"github.com/virtual-go/vos"
)

func NewEnvVarInputSrc() InputSrc {
	return envVarInputSrc{
		os:          vos.New(osfs.New()),
		readHistory: map[string]struct{}{},
	}
}

// envVarInputSrc implements InputSrc interface by sourcing inputs from env vars
type envVarInputSrc struct {
	os          vos.VOS
	readHistory map[string]struct{} // tracks reads
}

func (this envVarInputSrc) Read(
	inputName string,
) *string {
	if _, ok := this.readHistory[inputName]; ok {
		// enforce read at most once.
		return nil
	}

	if inputValue := this.os.Getenv(inputName); "" != inputValue {
		// track read history
		this.readHistory[inputName] = struct{}{}

		return &inputValue
	}
	return nil
}
