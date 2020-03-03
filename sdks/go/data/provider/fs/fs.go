package fs

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/sdks/go/data/provider"
	"github.com/opctl/opctl/sdks/go/model"
	"os"
	"path/filepath"
)

// New returns a data provider which sources pkgs from the filesystem
func New(
	basePaths ...string,
) provider.Provider {
	return _fs{
		os:        ios.New(),
		basePaths: basePaths,
	}
}

type _fs struct {
	os        ios.IOS
	basePaths []string
}

func (fp _fs) TryResolve(
	ctx context.Context,
	dataRef string,
) (model.DataHandle, error) {

	if filepath.IsAbs(dataRef) {
		if _, err := fp.os.Stat(dataRef); nil != err && !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}
		return newHandle(dataRef), nil
	}

	for _, basePath := range fp.basePaths {

		// attempt to resolve from basePath/.opspec dir
		testPath := filepath.Join(basePath, provider.DotOpspecDirName, dataRef)
		_, err := fp.os.Stat(testPath)
		if err == nil {
			return newHandle(testPath), nil
		}
		if nil != err && !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}

		// attempt to resolve from basePath
		testPath = filepath.Join(basePath, dataRef)
		_, err = fp.os.Stat(testPath)
		if err == nil {
			return newHandle(testPath), nil
		}
		if nil != err && !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}
	}

	return nil, nil
}
