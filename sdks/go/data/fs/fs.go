package fs

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
)

// New returns a data provider which sources data from the filesystem
func New() model.DataProvider {
	return _fs{}
}

type _fs struct{}

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
				return nil, fmt.Errorf("%w: path \"%s\"", model.ErrDataNotFoundResolution{}, dataRef)
			}
			return nil, err
		}
		return newHandle(dataRef), nil
	}

	return nil, fmt.Errorf("%w: path \"%s\"", model.ErrDataNotFoundResolution{}, dataRef)
}
