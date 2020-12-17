package cliparamsatisfier

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/ghodss/yaml"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params"
)

// CLIParamSatisfier attempts to satisfy the provided inputs via the provided inputSourcer
//
// if all fails an error is logged and we exit with a nonzero code.
//counterfeiter:generate -o fakes/cliParamSatisfier.go . CLIParamSatisfier
type CLIParamSatisfier interface {
	InputSrcFactory

	Satisfy(
		inputSourcer InputSourcer,
		inputs map[string]*model.Param,
	) (map[string]*model.Value, error)
}

func New(
	cliOutput clioutput.CliOutput,
) CLIParamSatisfier {

	return &_CLIParamSatisfier{
		cliOutput:       cliOutput,
		InputSrcFactory: newInputSrcFactory(),
	}
}

type _CLIParamSatisfier struct {
	cliOutput clioutput.CliOutput
	InputSrcFactory
}

func (cps _CLIParamSatisfier) Satisfy(
	inputSourcer InputSourcer,
	inputs map[string]*model.Param,
) (map[string]*model.Value, error) {

	argMap := map[string]*model.Value{}
	for _, paramName := range cps.getSortedParamNames(inputs) {
		param := inputs[paramName]

	paramLoop:
		for {
			var arg *model.Value

			rawArg, ok := inputSourcer.Source(paramName)
			if !ok {
				return nil, fmt.Errorf(`
-
  Prompt for "%v" failed; running in non-interactive terminal
-`, paramName)
			}

			switch {
			case nil == rawArg:
				// handle nil (returned by inputSourcer.Source for static defaults)
				break paramLoop
			case nil != param.Array:
				argValue := &[]interface{}{}
				argJSONBytes, err := yaml.YAMLToJSON([]byte(*rawArg))
				if nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				err = json.Unmarshal(argJSONBytes, argValue)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &model.Value{Array: argValue}
			case nil != param.Boolean:
				var err error
				if arg, err = coerce.ToBoolean(&model.Value{String: rawArg}); nil != err {
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
				arg = &model.Value{Dir: &absPath}
			case nil != param.File:
				absPath, err := filepath.Abs(*rawArg)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &model.Value{File: &absPath}
			case nil != param.Number:
				var err error
				if arg, err = coerce.ToNumber(&model.Value{String: rawArg}); nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
			case nil != param.Object:
				argValue := &map[string]interface{}{}
				argJSONBytes, err := yaml.YAMLToJSON([]byte(*rawArg))
				if nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				err = json.Unmarshal(argJSONBytes, argValue)
				if nil != err {
					// param not satisfied; notify & re-attempt!
					cps.notifyOfArgErrors([]error{err}, paramName)
					continue
				}
				arg = &model.Value{Object: argValue}
			case nil != param.Socket:
				arg = &model.Value{Socket: rawArg}
			case nil != param.String:
				arg = &model.Value{String: rawArg}
			}

			validateErr := params.Validate(
				map[string]*model.Value{paramName: arg},
				map[string]*model.Param{paramName: param},
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

	return argMap, nil
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
		messageBuffer.WriteString(
			fmt.Sprintf(`
	- %v`,
				validationError.Error(),
			),
		)
	}

	messageBuffer.WriteString(
		fmt.Sprintf(`
%v
-`,
			messageBuffer.String(),
		),
	)
	this.cliOutput.Error(messageBuffer.String())
}
