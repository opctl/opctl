package core

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/opctl/util/colorer"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/pkg/validate"
	"github.com/peterh/liner"
	"os"
	"strings"
)

// attempts to satisfy the provided params via:
// - checking the provided options
// - (falling back to) checking the environment
// - (falling back to) prompting for input
//
// if all fails an error is logged and we exit with a nonzero code.
type paramSatisfier interface {
	Satisfy(
		options []string,
		params []*model.Param,
	) (argMap map[string]*model.Data)
}

func newParamSatisfier(
	colorer colorer.Colorer,
	exiter exiter,
	validate validate.Validate,
) paramSatisfier {

	return &_paramSatisfier{
		colorer:  colorer,
		exiter:   exiter,
		validate: validate,
	}
}

type _paramSatisfier struct {
	colorer  colorer.Colorer
	exiter   exiter
	validate validate.Validate
}

func (this _paramSatisfier) Satisfy(
	options []string,
	params []*model.Param,
) (argMap map[string]*model.Data) {

	rawArgMap := make(map[string]string)
	for _, rawArg := range options {
		argParts := strings.Split(rawArg, "=")

		argName := argParts[0]
		var argValue string
		if len(argParts) > 1 {
			argValue = argParts[1]
		} else {
			argValue = os.Getenv(rawArg)
		}

		rawArgMap[argName] = argValue
	}

	argMap = make(map[string]*model.Data)
	paramIndex := 0
paramLoop:
	for paramIndex < len(params) {
		param := params[paramIndex]
		var paramName, rawArg string
		var arg *model.Data
		var argErrors []error

		switch {
		case nil != param.String:
			// obtain raw value
			stringParam := param.String
			paramName = stringParam.Name

			if providedArg, ok := rawArgMap[paramName]; ok {
				rawArg = providedArg
			} else if "" != os.Getenv(paramName) {
				rawArg = os.Getenv(paramName)
			} else if "" != stringParam.Default {
				// default value exists
				paramIndex++
				continue paramLoop
			} else {
				rawArg = this.promptForArg(paramName, stringParam.Description, stringParam.IsSecret)
			}
			arg = &model.Data{String: rawArg}
		case nil != param.Dir:
			// obtain raw value
			dirParam := param.Dir
			paramName = dirParam.Name

			if providedArg, ok := rawArgMap[paramName]; ok {
				rawArg = providedArg
			} else if "" != os.Getenv(paramName) {
				rawArg = os.Getenv(paramName)
			} else {
				rawArg = this.promptForArg(paramName, dirParam.Description, false)
			}
			arg = &model.Data{Dir: rawArg}
		case nil != param.File:
			// obtain raw value
			fileParam := param.File
			paramName = fileParam.Name

			if providedArg, ok := rawArgMap[paramName]; ok {
				rawArg = providedArg
			} else if "" != os.Getenv(paramName) {
				rawArg = os.Getenv(paramName)
			} else {
				rawArg = this.promptForArg(paramName, fileParam.Description, false)
			}
			arg = &model.Data{File: rawArg}
		case nil != param.Socket:
			socketParam := param.Socket
			paramName = socketParam.Name

			if providedArg, ok := rawArgMap[paramName]; ok {
				rawArg = providedArg
			} else if "" != os.Getenv(paramName) {
				rawArg = os.Getenv(paramName)
			} else {
				rawArg = this.promptForArg(paramName, socketParam.Description, false)
			}
			arg = &model.Data{Socket: rawArg}
		}

		// only perform semantic validation if no syntax errors exist
		// why? syntax errors imply semantic validation will fail
		if len(argErrors) < 1 {
			argErrors = append(argErrors, this.validate.Param(arg, param)...)
		}

		if len(argErrors) > 0 {
			this.notifyOfArgErrors(argErrors, paramName, rawArg)
			// we failed.. try again
			continue
		}

		// we succeeded.. store & move to next
		argMap[paramName] = arg
		paramIndex++
	}

	return

}

func (this _paramSatisfier) notifyOfArgErrors(
	errors []error,
	paramName string,
	rawArg string,
) {
	messageBuffer := bytes.NewBufferString(
		fmt.Sprintf(`
-
  %v invalid; provide valid value to proceed.
  Value: %v
  Error(s):`, paramName, rawArg),
	)
	for _, validationError := range errors {
		messageBuffer.WriteString(fmt.Sprintf(`
    - %v`, validationError.Error()))
	}
	fmt.Print(
		this.colorer.Error(`
%v
-`, messageBuffer.String()),
	)
}

func (this _paramSatisfier) promptForArg(
	paramName string,
	paramDescription string,
	paramIsSecret bool,
) (rawArg string) {

	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
	defer line.Close()

	fmt.Println(
		this.colorer.Attention(`
-
  Please provide value for parameter.
  Name: %v
  Description: %v
-`, paramName, paramDescription),
	)

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
		this.exiter.Exit(ExitReq{Message: message, Code: 1})
		return // support fake exiter
	}

	return
}
