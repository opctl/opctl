package cliparamsatisfier

import (
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc/cliprompt"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc/envvar"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc/paramdefault"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc/slice"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc/ymlfile"
	"github.com/opctl/opctl/sdks/go/model"
)

type InputSrcFactory interface {
	NewCliPromptInputSrc(
		inputs map[string]*model.Param,
	) inputsrc.InputSrc

	NewEnvVarInputSrc() inputsrc.InputSrc

	NewParamDefaultInputSrc(
		inputs map[string]*model.Param,
	) inputsrc.InputSrc

	NewSliceInputSrc(
		args []string,
		sep string,
	) inputsrc.InputSrc

	NewYMLFileInputSrc(
		filePath string,
	) (inputsrc.InputSrc, error)
}

func newInputSrcFactory() InputSrcFactory {
	return _inputSrcFactory{}
}

type _inputSrcFactory struct{}

func (is _inputSrcFactory) NewCliPromptInputSrc(
	inputs map[string]*model.Param,
) inputsrc.InputSrc {
	return cliprompt.New(inputs)
}

func (is _inputSrcFactory) NewEnvVarInputSrc() inputsrc.InputSrc {
	return envvar.New()
}

func (is _inputSrcFactory) NewParamDefaultInputSrc(
	inputs map[string]*model.Param,
) inputsrc.InputSrc {
	return paramdefault.New(inputs)
}

func (is _inputSrcFactory) NewSliceInputSrc(
	args []string,
	sep string,
) inputsrc.InputSrc {
	return slice.New(
		args,
		sep,
	)
}

func (is _inputSrcFactory) NewYMLFileInputSrc(
	filePath string,
) (inputsrc.InputSrc, error) {
	return ymlfile.New(filePath)
}
