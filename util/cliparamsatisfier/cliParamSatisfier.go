package cliparamsatisfier

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ CliParamSatisfier

import (
	"bytes"
	"fmt"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/clioutput"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall/inputs"
	"path/filepath"
	"sort"
	"strconv"
)

// attempts to satisfy the provided inputs via the provided inputSourcer
//
// if all fails an error is logged and we exit with a nonzero code.
type CliParamSatisfier interface {
	Satisfy(
		inputSourcer InputSourcer,
		inputs map[string]*model.Param,
	) map[string]*model.Data
}

func New(
	cliExiter cliexiter.CliExiter,
	cliOutput clioutput.CliOutput,
) CliParamSatisfier {

	return &_cliParamSatisfier{
		cliExiter: cliExiter,
		cliOutput: cliOutput,
		inputs:    inputs.New(),
	}
}

type _cliParamSatisfier struct {
	cliExiter cliexiter.CliExiter
	cliOutput clioutput.CliOutput
	inputs    inputs.Inputs
}

func (this _cliParamSatisfier) Satisfy(
	inputSourcer InputSourcer,
	inputs map[string]*model.Param,
) map[string]*model.Data {

	argMap := map[string]*model.Data{}
	for _, paramName := range this.getSortedParamNames(inputs) {
		param := inputs[paramName]

	paramLoop:
		for {
			var arg *model.Data

			rawArg := inputSourcer.Source(paramName)
			if nil == rawArg {
				msg := fmt.Sprintf(`
-
  Prompt for "%v" failed; running in non-interactive terminal
-`, paramName)
				this.cliExiter.Exit(cliexiter.ExitReq{Message: msg, Code: 1})
			}

			switch {
			case nil != param.String:
				arg = &model.Data{String: rawArg}
			case nil != param.Dir:
				absPath, err := filepath.Abs(*rawArg)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					this.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &model.Data{Dir: &absPath}
			case nil != param.File:
				absPath, err := filepath.Abs(*rawArg)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					this.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &model.Data{File: &absPath}
			case nil != param.Number:
				argValue, err := strconv.ParseFloat(*rawArg, 64)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					this.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &model.Data{Number: &argValue}
			case nil != param.Socket:
				arg = &model.Data{Socket: rawArg}
			}

			argErrors := this.inputs.Validate(
				map[string]*model.Data{paramName: arg},
				map[string]*model.Param{paramName: param},
			)
			if len(argErrors) > 0 {
				this.notifyOfArgErrors(argErrors[paramName], paramName)

				// param not satisfied; re-attempt it!
				continue
			}

			// param satisfied; store & move to next!
			argMap[paramName] = arg
			break paramLoop
		}
	}

	return argMap
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
) {
	messageBuffer := bytes.NewBufferString(
		fmt.Sprintf(`
-
  %v invalid; provide valid value to proceed.
  Error(s):`, paramName),
	)
	for _, validationError := range errors {
		messageBuffer.WriteString(fmt.Sprintf(`
    - %v`, validationError.Error()))
	}
	this.cliOutput.Error(`
%v
-`, messageBuffer.String())
}
