package cmd

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/array"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

// Interpret a container cmd
func Interpret(
	scope map[string]*ipld.Node,
	containerCallSpecCmd interface{},
) ([]string, error) {
	if containerCallSpecCmd == nil {
		return []string{}, nil
	}

	containerCallCmdArray, err := array.Interpret(
		scope,
		containerCallSpecCmd,
	)
	if err != nil {
		return nil, err
	}

	containerCallCmd := []string{}

	for _, cmdEntryExpression := range *containerCallCmdArray.Array {
		// interpret each entry as string
		cmdEntry, err := str.Interpret(scope, cmdEntryExpression)
		if err != nil {
			return nil, err
		}
		containerCallCmd = append(containerCallCmd, *cmdEntry.String)
	}

	return containerCallCmd, nil
}
