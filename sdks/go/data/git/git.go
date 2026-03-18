package git

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
	"golang.org/x/sync/singleflight"
)

// singleFlightGroup is used to ensure resolves don't race across provider intances
var resolveSingleFlightGroup singleflight.Group

// New returns a data provider which sources data from git repos
func New(
	basePath string,
	pullCreds *model.Creds,
) model.DataProvider {
	return _git{
		basePath:  basePath,
		pullCreds: pullCreds,
	}
}

type _git struct {
	basePath  string
	pullCreds *model.Creds
}

func (gp _git) Label() string {
	return "git"
}

func (gp _git) TryResolve(
	ctx context.Context,
	dataRef string,
) (model.DataHandle, error) {
	repoRef, err := parseRef(dataRef)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", model.ErrDataGitInvalidRef{}, err)
	}

	repoPath := repoRef.ToPath(gp.basePath)

	// attempt to resolve within singleFlight.Group to ensure concurrent resolves don't race
	if _, err, _ := resolveSingleFlightGroup.Do(
		repoPath,
		func() (interface{}, error) {
			completeMarkerPath := filepath.Join(gp.basePath, fmt.Sprintf(".%x", sha1.Sum([]byte(repoPath))))

			// lightweight check: resolve the tag to a commit SHA on the remote
			// without cloning (equivalent to git ls-remote)
			remoteHash, hashErr := resolveRemoteHash(ctx, repoRef, gp.pullCreds)
			if hashErr != nil {
				// can't reach remote — fall back to cache if available, else fail
				if _, err := os.Stat(completeMarkerPath); err == nil {
					return nil, nil
				}
				return nil, hashErr
			}

			// if the cached hash matches the remote, nothing has changed
			if cachedHash, err := os.ReadFile(completeMarkerPath); err == nil {
				if string(cachedHash) == remoteHash {
					return nil, nil
				}
			}

			// tag has moved (or no cache) — clone into a temp path so the
			// existing cache is never disturbed until the new clone succeeds
			tmpPath := repoPath + ".tmp"
			os.RemoveAll(tmpPath)

			if cloneErr := Clone(ctx, tmpPath, repoRef, gp.pullCreds); cloneErr != nil {
				// clone failed — fall back to cache if available
				os.RemoveAll(tmpPath)
				if _, err := os.Stat(completeMarkerPath); err == nil {
					return nil, nil
				}
				return nil, cloneErr
			}

			// atomically replace the cached copy and record the new hash
			os.RemoveAll(repoPath)
			if err := os.Rename(tmpPath, repoPath); err != nil {
				return nil, err
			}
			if err := os.WriteFile(completeMarkerPath, []byte(remoteHash), 0755); err != nil {
				return nil, err
			}
			return nil, nil
		},
	); err != nil {
		return nil, err
	}

	return newHandle(filepath.Join(gp.basePath, dataRef), dataRef), nil
}
