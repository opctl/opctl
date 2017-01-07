package core

//go:generate counterfeiter -o ./fakeWorkDirPathGetter.go --fake-name fakeWorkDirPathGetter ./ workDirPathGetter

import (
	"os"
	"path/filepath"
)

type workDirPathGetter interface {
	Get() (workDirPath string)
}

func newWorkDirPathGetter() workDirPathGetter {
	return _workDirPathGetter{}
}

type _workDirPathGetter struct{}

func (this _workDirPathGetter) Get() (workDirPath string) {
	wd, err := os.Getwd()
	workDirPath = filepath.ToSlash(wd)
	if err != nil {
		panic(err)
	}
	return
}
