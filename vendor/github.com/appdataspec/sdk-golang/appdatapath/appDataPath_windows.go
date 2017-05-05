package appdatapath

import (
	"fmt"
)

const (
	globalTemplate  = "%v"
	perUserTemplate = "%v"
)

func (this appDataPath) Global() (string, error) {
	programDataEnvVar := this.os.Getenv("PROGRAMDATA")
	if "" == programDataEnvVar {
		return "", errors.New("Unable to determine global app data path. Error was: PROGRAMDATA env var required")
	}
	return fmt.Sprintf(globalTemplate, programDataEnvVar), nil
}

func (this appDataPath) PerUser() (string, error) {
	localAppDataEnvVar := this.os.Getenv("LOCALAPPDATA")
	if "" == localAppDataEnvVar {
		return "", errors.New("Unable to determine per user app data path. Error was: LOCALAPPDATA env var required")
	}
	return fmt.Sprintf(perUserTemplate, localAppDataEnvVar), nil
}
