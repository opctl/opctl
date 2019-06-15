package inputs

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"bytes"
	"fmt"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/op/inputs/input"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/op/params"
)

type Interpreter interface {
	// Interpret interprets inputs via the provided inputArgs, inputParams, and scope;
	// opScratchDir will be used to store any run data such as type coercions to files
	Interpret(
		inputArgs map[string]interface{},
		inputParams map[string]*model.Param,
		parentOpHandle model.DataHandle,
		opPath string,
		scope map[string]*model.Value,
		opScratchDir string,
	) (map[string]*model.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		inputInterpreter: input.NewInterpreter(),
		paramsValidator:  params.NewValidator(),
		paramsDefaulter:  params.NewDefaulter(),
	}
}

type _interpreter struct {
	inputInterpreter input.Interpreter
	paramsValidator  params.Validator
	paramsDefaulter  params.Defaulter
}

func (itp _interpreter) Interpret(
	inputArgs map[string]interface{},
	inputParams map[string]*model.Param,
	parentOpHandle model.DataHandle,
	opPath string,
	scope map[string]*model.Value,
	opScratchDir string,
) (map[string]*model.Value, error) {
	interpretedArgs := map[string]*model.Value{}

	// 1) interpret
	paramErrMap := map[string]error{}
	for argName, argValue := range inputArgs {
		var err error
		interpretedArgs[argName], err = itp.inputInterpreter.Interpret(
			argName,
			argValue,
			inputParams[argName],
			parentOpHandle,
			scope,
			opScratchDir,
		)
		if nil != err {
			paramErrMap[argName] = err
		}
	}

	if len(paramErrMap) > 0 {
		// return error w/ fancy formatted msg
		messageBuffer := bytes.NewBufferString("")
		for paramName, err := range paramErrMap {
			messageBuffer.WriteString(fmt.Sprintf(`
		- %v: %v`, paramName, err.Error()))
		}
		messageBuffer.WriteString(`
`)
		return nil, fmt.Errorf(`
-
  validation error(s):
%v
-`, messageBuffer.String())
	}

	// 2) apply defaults
	argsWithDefaults := itp.paramsDefaulter.Default(interpretedArgs, inputParams, opPath)

	// 3) validate
	return argsWithDefaults, itp.paramsValidator.Validate(argsWithDefaults, inputParams)
}
