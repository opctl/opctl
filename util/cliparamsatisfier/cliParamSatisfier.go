package cliparamsatisfier

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ CliParamSatisfier

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/opctl/util/clicolorer"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/opctl/util/clioutput"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/validate"
	"github.com/peterh/liner"
	"path/filepath"
	"sort"
	"strconv"
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
	for _, paramName := range this.getSortedParamNames(params) {
		param := params[paramName]
		// track if using the env provided value has been attempted
		var (
			isDefaultAttempted  bool
			isEnvAttempted      bool
			isProvidedAttempted bool
		)

	paramLoop:
		for {
			var (
				arg                *model.Data
				argErrors          []error
				rawArgDisplayValue string
			)
			switch {
			case nil != param.String:
				// obtain raw value
				stringParam := param.String
				var rawArg string

				if stringParam.IsSecret {
					rawArgDisplayValue = "************"
				}

				if providedArg, ok := rawArgMap[paramName]; ok && !isProvidedAttempted {
					// provided & we've not made any attempt to use it
					isProvidedAttempted = true
					rawArg = providedArg
				} else if "" != this.vos.Getenv(paramName) && !isEnvAttempted {
					// env var exists & we've not made any attempt to use it
					isEnvAttempted = true
					rawArg = this.vos.Getenv(paramName)
				} else if "" != stringParam.Default && !isDefaultAttempted {
					isDefaultAttempted = true
					rawArg = stringParam.Default
				} else {
					rawArg = this.promptForArg(paramName, stringParam.Description, stringParam.IsSecret)
				}

				if "" == rawArgDisplayValue {
					rawArgDisplayValue = rawArg
				}
				arg = &model.Data{String: rawArg}
			case nil != param.Dir:
				// obtain raw value
				dirParam := param.Dir
				var rawArg string

				if providedArg, ok := rawArgMap[paramName]; ok && !isProvidedAttempted {
					// provided & we've not made any attempt to use it
					isProvidedAttempted = true
					rawArg = providedArg
				} else if "" != this.vos.Getenv(paramName) && !isEnvAttempted {
					// env var exists & we've not made any attempt to use it
					isEnvAttempted = true
					rawArg = this.vos.Getenv(paramName)
				} else if "" != dirParam.Default && !isDefaultAttempted {
					isDefaultAttempted = true
					wd, err := this.vos.Getwd()
					if nil != err {
						argErrors = append(argErrors, err)
					}
					rawArg = filepath.Join(wd, dirParam.Default)
				} else {
					rawArg = this.promptForArg(paramName, dirParam.Description, false)
				}

				rawArgDisplayValue = rawArg
				arg = &model.Data{Dir: rawArg}
			case nil != param.File:
				// obtain raw value
				fileParam := param.File
				var rawArg string

				if providedArg, ok := rawArgMap[paramName]; ok && !isProvidedAttempted {
					// provided & we've not made any attempt to use it
					isProvidedAttempted = true
					rawArg = providedArg
				} else if "" != this.vos.Getenv(paramName) && !isEnvAttempted {
					// env var exists & we've not made any attempt to use it
					isEnvAttempted = true
					rawArg = this.vos.Getenv(paramName)
				} else if "" != fileParam.Default && !isDefaultAttempted {
					isDefaultAttempted = true
					wd, err := this.vos.Getwd()
					if nil != err {
						argErrors = append(argErrors, err)
					}
					rawArg = filepath.Join(wd, fileParam.Default)
				} else {
					rawArg = this.promptForArg(paramName, fileParam.Description, false)
				}

				rawArgDisplayValue = rawArg
				arg = &model.Data{File: rawArg}
			case nil != param.Number:
				// obtain raw value
				numberParam := param.Number
				var rawArg float64

				if numberParam.IsSecret {
					rawArgDisplayValue = "************"
				}

				var err error
				if providedArg, ok := rawArgMap[paramName]; ok && !isProvidedAttempted {
					// provided & we've not made any attempt to use it
					isProvidedAttempted = true
					rawArg, err = strconv.ParseFloat(providedArg, 64)
				} else if "" != this.vos.Getenv(paramName) && !isEnvAttempted {
					// env var exists & we've not made any attempt to use it
					isEnvAttempted = true
					rawArg, err = strconv.ParseFloat(this.vos.Getenv(paramName), 64)
				} else if 0 != numberParam.Default && !isDefaultAttempted {
					isDefaultAttempted = true
					rawArg = numberParam.Default
				} else {
					rawArg, err =
						strconv.ParseFloat(this.promptForArg(paramName, numberParam.Description, numberParam.IsSecret), 64)
				}
				if nil != err {
					argErrors = append(argErrors, err)
				}

				if "" == rawArgDisplayValue {
					rawArgDisplayValue = strconv.FormatFloat(rawArg, 'f', -1, 64)
				}
				arg = &model.Data{Number: rawArg}
			case nil != param.Socket:
				socketParam := param.Socket
				var rawArg string

				if providedArg, ok := rawArgMap[paramName]; ok && !isProvidedAttempted {
					// provided & we've not made any attempt to use it
					isProvidedAttempted = true
					rawArg = providedArg
				} else if "" != this.vos.Getenv(paramName) && !isEnvAttempted {
					// env var exists & we've not made any attempt to use it
					isEnvAttempted = true
					rawArg = this.vos.Getenv(paramName)
				} else {
					rawArg = this.promptForArg(paramName, socketParam.Description, false)
				}

				rawArgDisplayValue = rawArg
				arg = &model.Data{Socket: rawArg}
			}

			if len(argErrors) == 0 {
				// only validate args if there wasn't an error retrieving them
				argErrors = append(argErrors, this.validate.Param(arg, param)...)
			}

			if len(argErrors) > 0 {
				this.notifyOfArgErrors(argErrors, paramName, rawArgDisplayValue)

				// param not satisfied; re-attempt it!
				continue
			}

			// param satisfied; store & move to next!
			argMap[paramName] = arg
			break paramLoop
		}
	}

	return
}

func (this _cliParamSatisfier) getSortedParamNames(
	params map[string]*model.Param,
) []string {
	paramNames := []string{}
	for paramname := range params {
		paramNames = append(paramNames, paramname)
	}
	sort.Strings(paramNames)
	return paramNames
}

func (this _cliParamSatisfier) notifyOfArgErrors(
	errors []error,
	paramName string,
	displayValue string,
) {
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

//@TODO: add promptForEnumArg

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
