package cmd

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	stringPkg "github.com/opctl/opctl/sdks/go/opspec/interpreter/string"
	"github.com/opctl/opctl/sdks/go/types"
)

type Interpreter interface {
	Interpret(
		scope map[string]*types.Value,
		scgContainerCallCmd []interface{},
		opHandle types.DataHandle,
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
	scope map[string]*types.Value,
	scgContainerCallCmd []interface{},
	opHandle types.DataHandle,
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
