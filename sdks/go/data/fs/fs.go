package fs

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	aggregateError "github.com/opctl/opctl/sdks/go/aggregate_error"
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
				return nil, fmt.Errorf("path %s not found", dataRef)
			}
			return nil, err
		}
		return newHandle(dataRef), nil
	}

	var aggregateErr aggregateError.ErrAggregate

	if len(fp.basePaths) == 0 {
		return nil, errors.New("skipped")
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
		aggregateErr.AddError(fmt.Errorf("path %s not found", testPath))
	}

	return nil, aggregateErr
}
