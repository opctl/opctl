package git

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/data/provider"
	"github.com/opctl/opctl/sdks/go/data/provider/fs"
	"github.com/opctl/opctl/sdks/go/model"
	"golang.org/x/sync/singleflight"
)

// singleFlightGroup is used to ensure resolves don't race across provider intances
var resolveSingleFlightGroup singleflight.Group

// New returns a data provider which sources pkgs from git repos
func New(
	basePath string,
	pullCreds *model.PullCreds,
) provider.Provider {
	return _git{
		localFSProvider: fs.New(basePath),
		basePath:        basePath,
		puller:          newPuller(),
		pullCreds:       pullCreds,
	}
}

type _git struct {
	// composed of fsProvider
	localFSProvider provider.Provider
	basePath        string
	puller          puller
	pullCreds       *model.PullCreds
}

func (gp _git) TryResolve(
	ctx context.Context,
	dataRef string,
) (model.DataHandle, error) {

	// attempt to resolve within singleFlight.Group to ensure concurrent resolves don't race
	handle, err, _ := resolveSingleFlightGroup.Do(
		dataRef,
		func() (interface{}, error) {
			// attempt to resolve from cache
			handle, err := gp.localFSProvider.TryResolve(ctx, dataRef)
			if nil != err {
				return nil, err
			} else if nil != handle {
				return handle, nil
			}

			// attempt pull if cache miss
			err = gp.puller.Pull(ctx, gp.basePath, dataRef, gp.pullCreds)
			if nil != err {
				return nil, err
			}
			return newHandle(filepath.Join(gp.basePath, dataRef), dataRef), nil
		},
	)

	if nil != err {
		return nil, err
	}
	return handle.(model.DataHandle), nil
}
