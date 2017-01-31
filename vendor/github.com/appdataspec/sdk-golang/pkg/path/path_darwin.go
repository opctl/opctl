package path

import (
	"fmt"
)

const (
	globalTemplate  = "/Library/Application Support"
	perUserTemplate = "%v/Library/Application Support"
)

func (this path) Global() string {
	return globalTemplate
}

func (this path) PerUser() string {
	localAppDataEnvVar := this.vos.Getenv("HOME")
	if "" == localAppDataEnvVar {
		panic("Unable to determine per user app data path. Error was: HOME env var required")
	}
	return fmt.Sprintf(
		perUserTemplate,
		localAppDataEnvVar,
	)
}
