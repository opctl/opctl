package git

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/opctl/opctl/sdks/go/model"
)

// Clone 'dataRef' to 'path'
// nil pullCreds will be ignored
//
// expected errs:
//   - ErrDataProviderAuthentication on authentication failure
//   - ErrDataProviderAuthorization on authorization failure
func Clone(
	ctx context.Context,
	repoPath string,
	repoRef *ref,
	authOpts *model.Creds,
) error {

	cloneOptions := &git.CloneOptions{
		URL:           fmt.Sprintf("https://%v", repoRef.Name),
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%v", repoRef.Version)),
		Depth:         1,
	}

	if authOpts != nil {
		cloneOptions.Auth = &http.BasicAuth{
			Username: authOpts.Username,
			Password: authOpts.Password,
		}
	}

	if _, err := git.PlainCloneContext(
		ctx,
		repoPath,
		false,
		cloneOptions,
	); err != nil {
		if _, ok := err.(git.NoMatchingRefSpecError); ok {
			return fmt.Errorf("%w: version \"%s\"", model.ErrDataRefResolution{}, repoRef.Version)
		}
		if errors.Is(err, transport.ErrAuthenticationRequired) {
			return model.ErrDataProviderAuthentication{}
		}
		if errors.Is(err, transport.ErrAuthorizationFailed) {
			return model.ErrDataProviderAuthorization{}
		}
		return err
	}

	return os.RemoveAll(filepath.Join(repoPath, ".git"))

}
