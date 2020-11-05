package fs

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"os"
	"path/filepath"

	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/sdks/go/model"
)

// New returns a data provider which sources pkgs from the filesystem
func New(
	basePaths ...string,
) model.DataProvider {
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
		_, err := fp.os.Stat(dataRef)
		if nil == err {
			return newHandle(dataRef), nil
		} else if !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}
		return nil, nil
	}

	for _, basePath := range fp.basePaths {

		// attempt to resolve from basePath
		testPath := filepath.Join(basePath, dataRef)
		_, err := fp.os.Stat(testPath)
		if nil == err {
			return newHandle(testPath), nil
		} else if !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}
	}

	return nil, nil
}
