package appdatapath

import (
	"fmt"
)

const (
	globalTemplate  = "%v"
	perUserTemplate = "%v"
)

func (this appDataPath) Global() string {
	programDataEnvVar := this.os.Getenv("PROGRAMDATA")
	if "" == programDataEnvVar {
		panic("Unable to determine global app data path. Error was: PROGRAMDATA env var required")
	}
	return fmt.Sprintf(globalTemplate, programDataEnvVar)
}

func (this appDataPath) PerUser() string {
	localAppDataEnvVar := this.os.Getenv("LOCALAPPDATA")
	if "" == localAppDataEnvVar {
		panic("Unable to determine per user app data path. Error was: LOCALAPPDATA env var required")
	}
	return fmt.Sprintf(
		perUserTemplate,
		localAppDataEnvVar,
	)
}
