package cmd

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

//Interpret a container cmd
func Interpret(
	scope map[string]*model.Value,
	containerCallSpecCmd []interface{},
) ([]string, error) {
	containerCallCmd := []string{}

	for _, cmdEntryExpression := range containerCallSpecCmd {
		// interpret each entry as string
		cmdEntry, err := str.Interpret(scope, cmdEntryExpression)
		if err != nil {
			return nil, err
		}
		containerCallCmd = append(containerCallCmd, *cmdEntry.String)
	}

	return containerCallCmd, nil
}
