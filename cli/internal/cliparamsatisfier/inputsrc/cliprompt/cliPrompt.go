package cliprompt

import (
	"fmt"
	"os"

	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/peterh/liner"
)

func New(
	inputs map[string]*model.Param,
) inputsrc.InputSrc {
	return cliPromptInputSrc{
		inputs:    inputs,
		cliOutput: clioutput.New(clicolorer.New(), os.Stderr, os.Stdout),
	}
}

// cliPromptInputSrc implements InputSrc interface by sourcing inputs from std in
type cliPromptInputSrc struct {
	inputs    map[string]*model.Param
	cliOutput clioutput.CliOutput
}

func (this cliPromptInputSrc) ReadString(
	inputName string,
) (*string, bool) {
	if param := this.inputs[inputName]; param != nil {
		var (
			isSecret    bool
			description string
		)

		switch {
		case param.Array != nil:
			isSecret = param.Array.IsSecret
			// @TODO remove after deprecation period
			description = param.Array.Description
		case param.Boolean != nil:
			// @TODO remove after deprecation period
			description = param.Boolean.Description
		case param.Dir != nil:
			// @TODO remove after deprecation period
			description = param.Dir.Description
		case param.File != nil:
			// @TODO remove after deprecation period
			description = param.File.Description
		case param.Number != nil:
			isSecret = param.Number.IsSecret
			// @TODO remove after deprecation period
			description = param.Number.Description
		case param.Object != nil:
			// @TODO remove after deprecation period
			description = param.Object.Description
		case param.Socket != nil:
			// @TODO remove after deprecation period
			description = param.Socket.Description
		case param.String != nil:
			isSecret = param.String.IsSecret
			// @TODO remove after deprecation period
			description = param.String.Description
		}

		if param.Description != "" {
			// non-deprecated property takes precedence
			description = param.Description
		}

		line := liner.NewLiner()
		defer line.Close()
		line.SetCtrlCAborts(true)

		this.cliOutput.Attention(
			fmt.Sprintf(`
-
  Please provide '%v'.
  Description: %v
-`,
				inputName,
				description,
			),
		)

		// liner has inconsistent behavior if non empty prompt arg passed so use ""
		var (
			err    error
			rawArg string
		)
		if isSecret {
			rawArg, err = line.PasswordPrompt("")
		} else {
			rawArg, err = line.Prompt("")
		}
		if err == nil {
			return &rawArg, true
		}
	}

	return nil, false
}
