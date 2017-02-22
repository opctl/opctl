package appdatapath

import (
	"fmt"
)

const (
	globalTemplate  = "/Library/Application Support"
	perUserTemplate = "%v/Library/Application Support"
)

func (this appDataPath) Global() string {
	return globalTemplate
}

func (this appDataPath) PerUser() string {
	localAppDataEnvVar := this.os.Getenv("HOME")
	if "" == localAppDataEnvVar {
		panic("Unable to determine per user app data path. Error was: HOME env var required")
	}
	return fmt.Sprintf(
		perUserTemplate,
		localAppDataEnvVar,
	)
}
