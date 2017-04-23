package pkg

import (
	"path"
)

// Resolve resolves a local package according to opspec package resolution rules and returns it's absolute path.
func (this pkg) Resolve(
	basePath,
	pkgRef string,
) (string, bool) {
	var testPath string

	// 1. attempt to resolve from basePath/.opspec dir
	testPath = path.Join(basePath, DotOpspecDirName, pkgRef)
	if _, err := this.os.Stat(testPath); nil == err {
		return testPath, true
	}

	// 2. attempt to resolve from basePath
	testPath = path.Join(basePath, pkgRef)
	if _, err := this.os.Stat(testPath); nil == err {
		return testPath, true
	}

	// 3. attempt to resolve from cache
	if testPath, err := constructCachePath(pkgRef); nil == err {
		if _, err = this.os.Stat(testPath); nil == err {
			return testPath, true
		}
	}

	// 4. giveup
	return "", false
}
