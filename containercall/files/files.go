package files

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Files

import (
	"github.com/golang-interfaces/iio"
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/filecopier"
	"github.com/opspec-io/sdk-golang/model"
)

type Files interface {
	Interpret(
		pkgPath string,
		scope map[string]*model.Value,
		scgContainerCallFiles map[string]string,
		scratchDirPath string,
	) (map[string]string, error)
}

func New(
	rootFSPath string,
) Files {
	return _Files{
		fileCopier: filecopier.New(),
		io:         iio.New(),
		os:         ios.New(),
		rootFSPath: rootFSPath,
	}
}

type _Files struct {
	fileCopier filecopier.FileCopier
	io         iio.IIO
	os         ios.IOS
	rootFSPath string
}
