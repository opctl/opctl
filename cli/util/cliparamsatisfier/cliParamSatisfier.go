package cliparamsatisfier

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ CLIParamSatisfier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/ghodss/yaml"
	"github.com/opctl/opctl/cli/util/cliexiter"
	"github.com/opctl/opctl/cli/util/clioutput"
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params"
	"github.com/opctl/opctl/sdks/go/types"
)

// CLIParamSatisfier attempts to satisfy the provided inputs via the provided inputSourcer
//
// if all fails an error is logged and we exit with a nonzero code.
type CLIParamSatisfier interface {
	InputSrcFactory

	Satisfy(
		inputSourcer InputSourcer,
		inputs map[string]*types.Param,
	) map[string]*types.Value
}

func New(
	cliExiter cliexiter.CliExiter,
	cliOutput clioutput.CliOutput,
) CLIParamSatisfier {

	return &_CLIParamSatisfier{
		cliExiter:       cliExiter,
		cliOutput:       cliOutput,
		coerce:          coerce.New(),
		paramsValidator: params.NewValidator(),
		InputSrcFactory: newInputSrcFactory(),
	}
}

type _CLIParamSatisfier struct {
	cliExiter       cliexiter.CliExiter
	cliOutput       clioutput.CliOutput
	coerce          coerce.Coerce
	paramsValidator params.Validator
	InputSrcFactory
}

func (cps _CLIParamSatisfier) Satisfy(
	inputSourcer InputSourcer,
	inputs map[string]*types.Param,
) map[string]*types.Value {

	argMap := map[string]*types.Value{}
	for _, paramName := range cps.getSortedParamNames(inputs) {
		param := inputs[paramName]

	paramLoop:
		for {
			var arg *types.Value

			rawArg, ok := inputSourcer.Source(paramName)
			if !ok {
				msg := fmt.Sprintf(`
-
  Prompt for "%v" failed; running in non-interactive terminal
-`, paramName)
				cps.cliExiter.Exit(cliexiter.ExitReq{Message: msg, Code: 1})
			}

			switch {
			case nil == rawArg:
				// handle nil (returned by inputSourcer.Source for static defaults)
				break paramLoop
			case nil != param.Array:
				argValue := &[]interface{}{}
				argJsonBytes, err := yaml.YAMLToJSON([]byte(*rawArg))
				if nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				err = json.Unmarshal(argJsonBytes, argValue)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &types.Value{Array: argValue}
			case nil != param.Boolean:
				var err error
				if arg, err = cps.coerce.ToBoolean(&types.Value{String: rawArg}); nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
			case nil != param.Dir:
				absPath, err := filepath.Abs(*rawArg)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &types.Value{Dir: &absPath}
			case nil != param.File:
				absPath, err := filepath.Abs(*rawArg)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &types.Value{File: &absPath}
			case nil != param.Number:
				var err error
				if arg, err = cps.coerce.ToNumber(&types.Value{String: rawArg}); nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
			case nil != param.Object:
				argValue := &map[string]interface{}{}
				argJsonBytes, err := yaml.YAMLToJSON([]byte(*rawArg))
				if nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				err = json.Unmarshal(argJsonBytes, argValue)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &types.Value{Object: argValue}
			case nil != param.Socket:
				arg = &types.Value{Socket: rawArg}
			case nil != param.String:
				arg = &types.Value{String: rawArg}
			}

			validateErr := cps.paramsValidator.Validate(
				map[string]*types.Value{paramName: arg},
				map[string]*types.Param{paramName: param},
			)
			if nil != validateErr {
				cps.notifyOfArgErrors([]error{validateErr}, paramName)

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
	params map[string]*types.Param,
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
