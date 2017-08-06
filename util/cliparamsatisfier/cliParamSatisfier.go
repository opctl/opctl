package cliparamsatisfier

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ CLIParamSatisfier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
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
type CLIParamSatisfier interface {
	InputSrcFactory

	Satisfy(
		inputSourcer InputSourcer,
		inputs map[string]*model.Param,
	) map[string]*model.Value
}

func New(
	cliExiter cliexiter.CliExiter,
	cliOutput clioutput.CliOutput,
) CLIParamSatisfier {

	return &_CLIParamSatisfier{
		cliExiter:       cliExiter,
		cliOutput:       cliOutput,
		inputs:          inputs.New(),
		InputSrcFactory: newInputSrcFactory(),
	}
}

type _CLIParamSatisfier struct {
	cliExiter cliexiter.CliExiter
	cliOutput clioutput.CliOutput
	inputs    inputs.Inputs
	InputSrcFactory
}

func (this _CLIParamSatisfier) Satisfy(
	inputSourcer InputSourcer,
	inputs map[string]*model.Param,
) map[string]*model.Value {

	argMap := map[string]*model.Value{}
	for _, paramName := range this.getSortedParamNames(inputs) {
		param := inputs[paramName]

	paramLoop:
		for {
			var arg *model.Value

			rawArg, ok := inputSourcer.Source(paramName)
			if !ok {
				msg := fmt.Sprintf(`
-
  Prompt for "%v" failed; running in non-interactive terminal
-`, paramName)
				this.cliExiter.Exit(cliexiter.ExitReq{Message: msg, Code: 1})
			}

			switch {
			case nil == rawArg:
				// handle nil (returned by inputSourcer.Source for static defaults)
				break paramLoop
			case nil != param.Dir:
				absPath, err := filepath.Abs(*rawArg)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					this.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &model.Value{Dir: &absPath}
			case nil != param.File:
				absPath, err := filepath.Abs(*rawArg)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					this.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &model.Value{File: &absPath}
			case nil != param.Number:
				argValue, err := strconv.ParseFloat(*rawArg, 64)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					this.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &model.Value{Number: &argValue}
			case nil != param.Object:
				argValue := map[string]interface{}{}
				argJsonBytes, err := yaml.YAMLToJSON([]byte(*rawArg))
				if nil != err {
					// param not satisfied; notify & re-attempt!
					this.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				err = json.Unmarshal(argJsonBytes, &argValue)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					this.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &model.Value{Object: argValue}
			case nil != param.Socket:
				arg = &model.Value{Socket: rawArg}
			case nil != param.String:
				arg = &model.Value{String: rawArg}
			}

			argErrors := this.inputs.Validate(
				map[string]*model.Value{paramName: arg},
				map[string]*model.Param{paramName: param},
			)
			if len(argErrors) > 0 {
				this.notifyOfArgErrors(argErrors[paramName], paramName)

				// param not satisfied; re-attempt it!
				continue
			}

			if nil != arg {
				// only store non-nil args
				argMap[paramName] = arg
			}

			// param satisfied; move to next!
			break paramLoop
		}
	}

	return argMap
}

func (this _CLIParamSatisfier) getSortedParamNames(
	params map[string]*model.Param,
) []string {
	paramNames := []string{}
	for paramname := range params {
		paramNames = append(paramNames, paramname)
	}
	sort.Strings(paramNames)
	return paramNames
}

func (this _CLIParamSatisfier) notifyOfArgErrors(
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
