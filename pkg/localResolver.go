package pkg

//go:generate counterfeiter -o ./fakeLocalResolver.go --fake-name fakeLocalResolver ./ localResolver

import (
	pathPkg "path"

	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	"github.com/virtual-go/vos"
)

type localResolver interface {
	Resolve(
		basePath string,
		pkgRef string,
	) (string, bool)
}

func newLocalResolver(
	os vos.VOS,
) localResolver {
	return _localResolver{
		os: os,
	}
}

type _localResolver struct {
	appDataSpec appdatapath.AppDataPath
	os          vos.VOS
}

func (this _localResolver) Resolve(
	basePath string,
	pkgRef string,
) (string, bool) {
	var testPath string

	// 1. attempt to resolve from basePath/.opspec dir
	testPath = pathPkg.Join(basePath, DotOpspecDirName, pkgRef)
	if _, err := this.os.Stat(testPath); nil == err {
		return testPath, true
	}

	// 2. attempt to resolve from basePath
	testPath = pathPkg.Join(basePath, pkgRef)
	if _, err := this.os.Stat(testPath); nil == err {
		return testPath, true
	}

	// 3. attempt to resolve from cache
	testPath = pathPkg.Join(
		appdatapath.New().PerUser(),
		"opspec",
		"cache",
		"pkgs",
		pkgRef,
	)
	if _, err := this.os.Stat(testPath); nil == err {
		return testPath, true
	}

	// 4. giveup
	return "", false
}
