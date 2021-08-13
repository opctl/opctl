package clioutput

import (
	"os"
	"path"
	"strings"
)

// FormatOpRef gives a more appropriate description of an op's reference
// Local ops will be formatted as paths relative to the working directory or
// home directory
func FormatOpRef(opRef string) string {
	if path.IsAbs(opRef) {
		cwd, err := os.Getwd()
		if err != nil {
			return opRef
		}

		if strings.HasPrefix(opRef, cwd) {
			return "." + opRef[len(cwd):]
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return opRef
		}

		if strings.HasPrefix(opRef, home) {
			return "~" + opRef[len(home):]
		}
	}
	return opRef
}
