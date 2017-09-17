package appdatapath

import (
	"errors"
	"fmt"
)

const (
	globalTemplate  = "/var/lib"
	perUserTemplate = "%v"
)

func (this appDataPath) Global() (string, error) {
	return globalTemplate, nil
}

func (this appDataPath) PerUser() (string, error) {
	localAppDataEnvVar := this.os.Getenv("HOME")
	if "" == localAppDataEnvVar {
		return "", errors.New("unable to determine per user app data path. Error was: HOME env var required")
	}
	return fmt.Sprintf(perUserTemplate, localAppDataEnvVar), nil
}
