package pkg

import (
	"path"
)

// Resolve attempts to resolve a package from lookPaths according to opspec package resolution rules.
// each lookPath will be tried in the provided order until either resolution succeeds or all lookPaths are exhausted
// if successful, the absolute path of the resolved package will be returned along w/ true
func (this pkg) Resolve(
	pkgRef string,
	lookPaths ...string,
) (string, bool) {

	for _, lookPath := range lookPaths {
		// 1. attempt to resolve from lookPath/.opspec dir
		testPath := path.Join(lookPath, DotOpspecDirName, pkgRef)
		if _, err := this.os.Stat(testPath); nil == err {
			return testPath, true
		}

		// 2. attempt to resolve from lookPath
		testPath = path.Join(lookPath, pkgRef)
		if _, err := this.os.Stat(testPath); nil == err {
			return testPath, true
		}
	}

	// resolution unsuccessful
	return "", false
}
