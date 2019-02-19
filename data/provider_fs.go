package data

import (
	"context"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/model"
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

func (fp fsProvider) TryResolve(
	ctx context.Context,
	dataRef string,
) (model.DataHandle, error) {

	if filepath.IsAbs(dataRef) {
		if _, err := fp.os.Stat(dataRef); nil != err && !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}
		return newFSHandle(dataRef), nil
	}

	for _, basePath := range fp.basePaths {

		// attempt to resolve from basePath/.opspec dir
		testPath := filepath.Join(basePath, DotOpspecDirName, dataRef)
		_, err := fp.os.Stat(testPath)
		if err == nil {
			return newFSHandle(testPath), nil
		}
		if nil != err && !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}

		// attempt to resolve from basePath
		testPath = filepath.Join(basePath, dataRef)
		_, err = fp.os.Stat(testPath)
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
