package pkg

//go:generate counterfeiter -o ./fakeResolver.go --fake-name fakeResolver ./ resolver

import (
	pathPkg "path"

	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	"github.com/virtual-go/vos"
)

type resolver interface {
	// Resolve resolves a local package according to opspec package resolution rules and returns it's absolute path.
	Resolve(
		basePath,
		pkgRef string,
	) (string, bool)
}

func newResolver(
	os vos.VOS,
) resolver {
	return _resolver{
		os: os,
	}
}

type _resolver struct {
	appDataSpec appdatapath.AppDataPath
	os          vos.VOS
}

func (this _resolver) Resolve(
	basePath,
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
