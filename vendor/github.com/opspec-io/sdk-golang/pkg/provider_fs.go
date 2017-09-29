package pkg

import (
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"os"
	"path/filepath"
)

func (pf _providerFactory) NewFSProvider(
	basePaths ...string,
) Provider {
	return fsProvider{
		os:        ios.New(),
		basePaths: basePaths,
	}
}

type fsProvider struct {
	os        ios.IOS
	basePaths []string
}

func (lfsp fsProvider) TryResolve(
	pkgRef string,
) (model.PkgHandle, error) {

	if filepath.IsAbs(pkgRef) {
		if _, err := lfsp.os.Stat(pkgRef); nil != err && !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}
		return newFSHandle(pkgRef), nil
	}

	for _, basePath := range lfsp.basePaths {

		// attempt to resolve from basePath/.opspec dir
		testPath := filepath.Join(basePath, DotOpspecDirName, pkgRef)
		_, err := lfsp.os.Stat(testPath)
		if err == nil {
			return newFSHandle(testPath), nil
		}
		if nil != err && !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}

		// attempt to resolve from basePath
		testPath = filepath.Join(basePath, pkgRef)
		_, err = lfsp.os.Stat(testPath)
		if err == nil {
			return newFSHandle(testPath), nil
		}
		if nil != err && !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}
	}

	return nil, nil
}
