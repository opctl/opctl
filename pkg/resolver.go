package pkg

import (
	"github.com/golang-interfaces/ios"
	"path/filepath"
)

type resolver interface {
	// Resolve attempts to resolve a package from lookPaths according to opspec package resolution rules.
	// each lookPath will be tried in the provided order until either resolution succeeds or all lookPaths are exhausted
	// if successful, the absolute path of the resolved package will be returned along w/ true
	Resolve(
		pkgRef *PkgRef,
		lookPaths ...string,
	) (
		string,
		bool,
	)
}

func newResolver() resolver {
	return _resolver{
		os: ios.New(),
	}
}

type _resolver struct {
	os ios.IOS
}

func (this _resolver) Resolve(
	pkgRef *PkgRef,
	lookPaths ...string,
) (
	string,
	bool,
) {
	for _, lookPath := range lookPaths {
		// 1. attempt to resolve from lookPath/.opspec dir
		testPath := pkgRef.ToPath(filepath.Join(lookPath, DotOpspecDirName))
		if _, err := this.os.Stat(testPath); nil == err {
			return testPath, true
		}

		// 2. attempt to resolve from lookPath
		testPath = pkgRef.ToPath(lookPath)
		if _, err := this.os.Stat(testPath); nil == err {
			return testPath, true
		}
	}

	// resolution unsuccessful
	return "", false
}
