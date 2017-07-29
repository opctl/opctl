package pkg

import (
	"github.com/golang-interfaces/ios"
	"os"
	"path/filepath"
)

// NewLocalFSProvider returns a pkg provider which sources pkgs from the local filesystem
func (pf _ProviderFactory) NewLocalFSProvider(
	basePaths ...string,
) Provider {
	return localFSProvider{
		os:        ios.New(),
		basePaths: basePaths,
	}
}

type localFSProvider struct {
	os        ios.IOS
	basePaths []string
}

func (lfsp localFSProvider) TryResolve(
	pkgRef string,
) (Handle, error) {

	if filepath.IsAbs(pkgRef) {
		if _, err := lfsp.os.Stat(pkgRef); nil != err && !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}
		return newLocalHandle(pkgRef), nil
	}

	for _, basePath := range lfsp.basePaths {

		// attempt to resolve from basePath/.opspec dir
		testPath := filepath.Join(basePath, DotOpspecDirName, pkgRef)
		_, err := lfsp.os.Stat(testPath)
		if err == nil {
			return newLocalHandle(testPath), nil
		}
		if nil != err && !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}

		// attempt to resolve from basePath
		testPath = filepath.Join(basePath, pkgRef)
		_, err = lfsp.os.Stat(testPath)
		if err == nil {
			return newLocalHandle(testPath), nil
		}
		if nil != err && !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}
	}

	return nil, nil
}
