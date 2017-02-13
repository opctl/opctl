package cliparamsatisfier

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ CliParamSatisfier

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/opctl/util/clicolorer"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/opctl/util/clioutput"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/pkg/validate"
	"github.com/peterh/liner"
	"strings"
)

// attempts to satisfy the provided params via:
// - checking the provided options
// - (falling back to) checking the environment
// - (falling back to) prompting for input
//
// if all fails an error is logged and we exit with a nonzero code.
type CliParamSatisfier interface {
	Satisfy(
		options []string,
		params map[string]*model.Param,
	) (argMap map[string]*model.Data)
}

func New(
	cliColorer clicolorer.CliColorer,
	cliExiter cliexiter.CliExiter,
	cliOutput clioutput.CliOutput,
	validate validate.Validate,
	vos vos.Vos,
) CliParamSatisfier {

	return &_cliParamSatisfier{
		cliColorer: cliColorer,
		cliExiter:  cliExiter,
		cliOutput:  cliOutput,
		validate:   validate,
		vos:        vos,
	}
}

type _cliParamSatisfier struct {
	cliColorer clicolorer.CliColorer
	cliExiter  cliexiter.CliExiter
	cliOutput  clioutput.CliOutput
	validate   validate.Validate
	vos        vos.Vos
}

func (this _cliParamSatisfier) Satisfy(
	options []string,
	params map[string]*model.Param,
) (argMap map[string]*model.Data) {

	rawArgMap := make(map[string]string)
	for _, rawArg := range options {
		argParts := strings.Split(rawArg, "=")

		argName := argParts[0]
		var argValue string
		if len(argParts) > 1 {
			argValue = argParts[1]
		} else {
			argValue = this.vos.Getenv(rawArg)
		}

		rawArgMap[argName] = argValue
	}

	argMap = make(map[string]*model.Data)
	for paramName, param := range params {
	paramLoop:
		for {
			var rawArg string
			var arg *model.Data
			var argErrors []error
			var hideValue bool
			switch {
			case nil != param.String:
				// obtain raw value
				stringParam := param.String
				hideValue = stringParam.IsSecret

				if providedArg, ok := rawArgMap[paramName]; ok {
					rawArg = providedArg
				} else if "" != this.vos.Getenv(paramName) {
					rawArg = this.vos.Getenv(paramName)
				} else if "" != stringParam.Default {
					// default value exists
					break paramLoop
				} else {
					rawArg = this.promptForArg(paramName, stringParam.Description, stringParam.IsSecret)
				}
				arg = &model.Data{String: rawArg}
			case nil != param.Dir:
				// obtain raw value
				dirParam := param.Dir

				if providedArg, ok := rawArgMap[paramName]; ok {
					rawArg = providedArg
				} else if "" != this.vos.Getenv(paramName) {
					rawArg = this.vos.Getenv(paramName)
				} else {
					rawArg = this.promptForArg(paramName, dirParam.Description, false)
				}
				arg = &model.Data{Dir: rawArg}
			case nil != param.File:
				// obtain raw value
				fileParam := param.File

				if providedArg, ok := rawArgMap[paramName]; ok {
					rawArg = providedArg
				} else if "" != this.vos.Getenv(paramName) {
					rawArg = this.vos.Getenv(paramName)
				} else {
					rawArg = this.promptForArg(paramName, fileParam.Description, false)
				}
				arg = &model.Data{File: rawArg}
			case nil != param.Socket:
				socketParam := param.Socket

				if providedArg, ok := rawArgMap[paramName]; ok {
					rawArg = providedArg
				} else if "" != this.vos.Getenv(paramName) {
					rawArg = this.vos.Getenv(paramName)
				} else {
					rawArg = this.promptForArg(paramName, socketParam.Description, false)
				}
				arg = &model.Data{Socket: rawArg}
			}

			// validate
			argErrors = append(argErrors, this.validate.Param(arg, param)...)

			if len(argErrors) > 0 {
				this.notifyOfArgErrors(argErrors, paramName, hideValue, rawArg)
				// we failed.. try again
				continue
			}

			// we succeeded.. store & move to next
			argMap[paramName] = arg
			break paramLoop
		}
	}

	return
}

func (this _cliParamSatisfier) notifyOfArgErrors(
	errors []error,
	paramName string,
	hideValue bool,
	rawArg string,
) {
	var displayValue string
	if hideValue {
		displayValue = "**********"
	} else {
		displayValue = rawArg
	}
	messageBuffer := bytes.NewBufferString(
		fmt.Sprintf(`
-
  %v invalid; provide valid value to proceed.
  Value: %v
  Error(s):`, paramName, displayValue),
	)
	for _, validationError := range errors {
		messageBuffer.WriteString(fmt.Sprintf(`
    - %v`, validationError.Error()))
	}
	this.cliOutput.Error(`
%v
-`, messageBuffer.String())
}

func (this _cliParamSatisfier) promptForArg(
	paramName string,
	paramDescription string,
	paramIsSecret bool,
) (rawArg string) {

	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
	defer line.Close()

	this.cliOutput.Attention(`
-
  Please provide value for parameter.
  Name: %v
  Description: %v
-`, paramName, paramDescription)

	// liner has inconsistent behavior if non empty prompt arg passed so use ""
	var err error
	if paramIsSecret {
		rawArg, err = line.PasswordPrompt("")
	} else {
		rawArg, err = line.Prompt("")
	}

	if nil != err {
		message := fmt.Sprintf(`
-
  Prompt for input parameter "%v" failed.

  To specify the parameter either:
    a) provide it explicitly to the run command (via the -a option)
    b) set it as an environment variable
    c) run the op from an interactive shell and enter it when prompted
-`,
			paramName,
		)

		// explicitly call line.Close(); os.Exit() doesn't allow defer'd statements to be observed
		line.Close()

		this.cliExiter.Exit(cliexiter.ExitReq{Message: message, Code: 1})
		return // support fake exiter
	}

	return
}
