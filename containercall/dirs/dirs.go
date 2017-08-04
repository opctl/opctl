package dirs

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Dirs

import (
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/dircopier"
	"github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type Dirs interface {
	Interpret(
		pkgHandle model.PkgHandle,
		scope map[string]*model.Value,
		scgContainerCallFiles map[string]string,
		scratchDirPath string,
	) (map[string]string, error)
}

func New(
	rootFSPath string,
) Dirs {
	return _Dirs{
		dirCopier:    dircopier.New(),
		interpolater: interpolater.New(),
		os:           ios.New(),
		rootFSPath:   rootFSPath,
	}
}

type _Dirs struct {
	dirCopier    dircopier.DirCopier
	interpolater interpolater.Interpolater
	os           ios.IOS
	rootFSPath   string
}
