package appdata

import (
	"fmt"
	"os"
)

const (
	globalPathTemplate  = "%v"
	perUserPathTemplate = "%v"
)

func (this appData) GlobalPath() string {
	programDataEnvVar := os.Getenv("PROGRAMDATA")
	if "" == programDataEnvVar {
		panic("Unable to determine global app data path. Error was: PROGRAMDATA env var required")
	}
	return fmt.Sprintf(globalPathTemplate, programDataEnvVar)
}

func (this appData) PerUserPath() string {
	localAppDataEnvVar := os.Getenv("LOCALAPPDATA")
	if "" == localAppDataEnvVar {
		panic("Unable to determine per user app data path. Error was: LOCALAPPDATA env var required")
	}
	return fmt.Sprintf(
		perUserPathTemplate,
		localAppDataEnvVar,
	)
}
