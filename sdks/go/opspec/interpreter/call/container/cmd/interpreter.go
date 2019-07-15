package cmd

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"github.com/opctl/sdk-golang/model"
	stringPkg "github.com/opctl/sdk-golang/opspec/interpreter/string"
)

type Interpreter interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallCmd []interface{},
		opHandle model.DataHandle,
	) ([]string, error)
}

// NewInterpreter returns a new Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		stringInterpreter: stringPkg.NewInterpreter(),
	}
}

type _interpreter struct {
	stringInterpreter stringPkg.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	scgContainerCallCmd []interface{},
	opHandle model.DataHandle,
) ([]string, error) {
	dcgContainerCallCmd := []string{}

	for _, cmdEntryExpression := range scgContainerCallCmd {
		// interpret each entry as string
		cmdEntry, err := itp.stringInterpreter.Interpret(scope, cmdEntryExpression, opHandle)
		if nil != err {
			return nil, err
		}
		dcgContainerCallCmd = append(dcgContainerCallCmd, *cmdEntry.String)
	}

	return dcgContainerCallCmd, nil
}
