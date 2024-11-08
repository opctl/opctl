package datadir

import (
	"fmt"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/internal/unsudo"
)

// DataDir is an interface exposing the functionality we require in conjunction with our "data dir".
type DataDir interface {
	// Path resolves the data dir path
	Path() string
}

// New returns an initialized data dir instance
func New(
	dataDirPath string,
) (DataDir, error) {
	resolvedPath, err := filepath.Abs(dataDirPath)
	if err != nil {
		return nil, fmt.Errorf("error initializing opctl data dir: %w", err)
	}

	if err := unsudo.CreateDir(resolvedPath); err != nil {
		return nil, fmt.Errorf("error initializing opctl data dir: %w", err)
	}

	return _datadir{
		resolvedPath: resolvedPath,
	}, err
}

type _datadir struct {
	resolvedPath string
}

func (dd _datadir) Path() string {
	return filepath.Join(dd.resolvedPath)
}
