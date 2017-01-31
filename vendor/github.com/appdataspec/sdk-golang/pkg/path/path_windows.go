package path

import (
	"fmt"
)

const (
	globalTemplate  = "%v"
	perUserTemplate = "%v"
)

func (this path) Global() string {
	programDataEnvVar := this.vos.Getenv("PROGRAMDATA")
	if "" == programDataEnvVar {
		panic("Unable to determine global app data path. Error was: PROGRAMDATA env var required")
	}
	return fmt.Sprintf(globalTemplate, programDataEnvVar)
}

func (this path) PerUser() string {
	localAppDataEnvVar := this.vos.Getenv("LOCALAPPDATA")
	if "" == localAppDataEnvVar {
		panic("Unable to determine per user app data path. Error was: LOCALAPPDATA env var required")
	}
	return fmt.Sprintf(
		perUserTemplate,
		localAppDataEnvVar,
	)
}
