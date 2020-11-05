package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/opctl/opctl/sdks/go/model"
)

// Pull pulls 'dataRef' to 'path'
// nil pullCreds will be ignored
//
// expected errs:
//  - ErrDataProviderAuthentication on authentication failure
//  - ErrDataProviderAuthorization on authorization failure
func Pull(
	ctx context.Context,
	path string,
	dataRef string,
	authOpts *model.Creds,
) error {

	parsedPkgRef, err := parseRef(dataRef)
	if nil != err {
		return err
	}

	opPath := parsedPkgRef.ToPath(path)

	cloneOptions := &git.CloneOptions{
		URL:           fmt.Sprintf("https://%v", parsedPkgRef.Name),
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%v", parsedPkgRef.Version)),
		Depth:         1,
		Progress:      os.Stdout,
	}

	if nil != authOpts {
		cloneOptions.Auth = &http.BasicAuth{
			Username: authOpts.Username,
			Password: authOpts.Password,
		}
	}

	if _, err := git.PlainClone(
		opPath,
		false,
		cloneOptions,
	); nil != err {
		switch err.Error() {
		case transport.ErrAuthenticationRequired.Error():
			return model.ErrDataProviderAuthentication{}
		case transport.ErrAuthorizationFailed.Error():
			return model.ErrDataProviderAuthorization{}
		case git.ErrRepositoryAlreadyExists.Error():
			return nil
			// NoOp on repo already exists
		default:
			return err
		}
	}

	// remove pkg '.git' sub dir
	return os.RemoveAll(filepath.Join(opPath, ".git"))

}
