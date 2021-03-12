package fs

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
)

// New returns a data provider which sources pkgs from the filesystem
func New(
	basePaths ...string,
) model.DataProvider {
	return _fs{
		basePaths: basePaths,
	}
}

type _fs struct {
	basePaths []string
}

func (fp _fs) Label() string {
	return "filesystem"
}

func (fp _fs) TryResolve(
	ctx context.Context,
	dataRef string,
) (model.DataHandle, error) {

	if filepath.IsAbs(dataRef) {
		if _, err := os.Stat(dataRef); err != nil {
			if os.IsNotExist(err) {
				return nil, errors.New("not found")
			}
			return nil, err
		}
		return newHandle(dataRef), nil
	}

	for _, basePath := range fp.basePaths {
		// attempt to resolve from basePath
		testPath := filepath.Join(basePath, dataRef)
		if _, err := os.Stat(testPath); err == nil {
			return newHandle(testPath), nil
		} else if !os.IsNotExist(err) {
			// don't return not found errors, instead continue checking other base paths
			return nil, err
		}
	}

	return nil, errors.New("not found")
}
