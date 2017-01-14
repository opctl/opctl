package appdata

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
)

//go:generate counterfeiter -o ./fakeAppData.go --fake-name FakeAppData ./ AppData

type AppData interface {
	GlobalPath() string
	PerUserPath() string
}

func New() AppData {
	return appData{}
}

type appData struct{}

func (this appData) getUserHomePath() string {
	userHomePath, err := homedir.Dir()
	if nil != err {
		panic(fmt.Sprintf("Unable to retrieve user home dir. Error was %v", err.Error()))
	}
	return userHomePath
}
