package cmd

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallCmd []interface{},
	) ([]string, error)
}

// NewInterpreter returns a new Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		stringInterpreter: str.NewInterpreter(),
	}
}

type _interpreter struct {
	stringInterpreter str.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	scgContainerCallCmd []interface{},
) ([]string, error) {
	dcgContainerCallCmd := []string{}

	for _, cmdEntryExpression := range scgContainerCallCmd {
		// interpret each entry as string
		cmdEntry, err := itp.stringInterpreter.Interpret(scope, cmdEntryExpression)
		if nil != err {
			return nil, err
		}
		dcgContainerCallCmd = append(dcgContainerCallCmd, *cmdEntry.String)
	}

	return dcgContainerCallCmd, nil
}
