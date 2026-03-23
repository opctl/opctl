package git

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
  "context"
  "crypto/sha1"
  "errors"
  "fmt"
  "io"
  "os"
  "path/filepath"

  "github.com/go-git/go-git/v5/plumbing/transport"
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
  // basePath is the local directory under which cloned git repos are cached
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

  repoAbsPath := repoRef.ToPath(gp.basePath)
  repoRelPath, _ := filepath.Rel(gp.basePath, repoAbsPath)
  tmpRelPath := repoRelPath + ".tmp"

  // attempt to resolve within singleFlight.Group to ensure concurrent resolves don't race
  if _, err, _ := resolveSingleFlightGroup.Do(
    repoAbsPath,
    func() (interface{}, error) {
      root, err := os.OpenRoot(gp.basePath)
      if err != nil {
        return nil, err
      }
      defer root.Close()

      markerName := fmt.Sprintf(".%x", sha1.Sum([]byte(repoRelPath)))

      // lightweight check: resolve the tag to a commit SHA on the remote
      // without cloning (equivalent to git ls-remote)
      remoteHash, hashErr := resolveRemoteHash(ctx, repoRef, gp.pullCreds)
      if hashErr != nil {
        if errors.Is(hashErr, transport.ErrAuthenticationRequired) {
          return nil, model.ErrDataProviderAuthentication{}
        }
        if errors.Is(hashErr, transport.ErrAuthorizationFailed) {
          return nil, model.ErrDataProviderAuthorization{}
        }
        return nil, hashErr
      }

      // if the cached hash matches the remote, nothing has changed
      if f, err := root.Open(markerName); err == nil {
        cachedHash, _ := io.ReadAll(f)
        f.Close()
        if string(cachedHash) == remoteHash {
          return nil, nil
        }
      }

      // tag has moved (or no cache) — clone into a temp path so the
      // existing cache is never disturbed until the new clone succeeds
      root.RemoveAll(tmpRelPath)

      if cloneErr := Clone(ctx, filepath.Join(gp.basePath, tmpRelPath), repoRef, gp.pullCreds); cloneErr != nil {
        root.RemoveAll(tmpRelPath)
        return nil, cloneErr
      }

      // atomically replace the cached copy and record the new hash
      root.RemoveAll(repoRelPath)
      if err := root.Rename(tmpRelPath, repoRelPath); err != nil {
        return nil, err
      }
      if err := root.WriteFile(markerName, []byte(remoteHash), 0600); err != nil {
        return nil, err
      }
      return nil, nil
    },
  ); err != nil {
    return nil, err
  }

  return newHandle(filepath.Join(gp.basePath, dataRef), dataRef), nil
}
