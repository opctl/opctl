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
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/sdks/go/model"
)

//counterfeiter:generate -o internal/fakes/puller.go . puller
type puller interface {
	// Pull pulls 'dataRef' to 'path'
	// nil pullCreds will be ignored
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	Pull(
		ctx context.Context,
		path string,
		dataRef string,
		pullCreds *model.PullCreds,
	) error
}

func newPuller() puller {
	return _puller{
		os:        ios.New(),
		refParser: newRefParser(),
	}
}

type _puller struct {
	os        ios.IOS
	refParser refParser
}

func (plr _puller) Pull(
	ctx context.Context,
	path string,
	dataRef string,
	authOpts *model.PullCreds,
) error {

	parsedPkgRef, err := plr.refParser.Parse(dataRef)
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
	return plr.os.RemoveAll(filepath.Join(opPath, ".git"))

}
