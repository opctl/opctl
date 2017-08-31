package pkg

//go:generate counterfeiter -o ./fakePuller.go --fake-name fakePuller ./ puller

import (
	"fmt"
	"github.com/golang-interfaces/gopkg.in-src-d-go-git.v4"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"golang.org/x/sync/singleflight"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
	"path/filepath"
)

// pullerSingleFlightGroup is used to ensures pulls don't race across puller intances
var pullerSingleFlightGroup singleflight.Group

type puller interface {
	// Pull pulls 'pkgRef' to 'path'
	// nil pullCreds will be ignored
	//
	// expected errs:
	//  - ErrPkgPullAuthentication on authentication failure
	//  - ErrPkgPullAuthorization on authorization failure
	Pull(
		path string,
		pkgRef string,
		pullCreds *model.PullCreds,
	) error
}

func newPuller() puller {
	return _puller{
		git:       igit.New(),
		os:        ios.New(),
		refParser: newRefParser(),
	}
}

type _puller struct {
	git       igit.IGit
	os        ios.IOS
	refParser refParser
}

func (this _puller) Pull(
	path string,
	pkgRef string,
	authOpts *model.PullCreds,
) error {

	parsedPkgRef, err := this.refParser.Parse(pkgRef)
	if nil != err {
		return err
	}

	pkgPath := parsedPkgRef.ToPath(path)

	// ensure only a single pull done at once
	_, err, _ = pullerSingleFlightGroup.Do(pkgPath, func() (interface{}, error) {
		cloneOptions := &git.CloneOptions{
			URL:           fmt.Sprintf("https://%v", parsedPkgRef.Name),
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%v", parsedPkgRef.Version)),
			// @TODO re-enable once https://github.com/src-d/go-git/issues/529 released
			// Depth:         1,
			Progress: os.Stdout,
		}

		if nil != authOpts {
			cloneOptions.Auth = http.NewBasicAuth(authOpts.Username, authOpts.Password)
		}

		if _, err := this.git.PlainClone(
			pkgPath,
			false,
			cloneOptions,
		); nil != err {
			switch err.Error() {
			case transport.ErrAuthenticationRequired.Error():
				// clone failed; cleanup remnants
				this.os.RemoveAll(pkgPath)
				return nil, model.ErrPkgPullAuthentication{}
			case transport.ErrAuthorizationFailed.Error():
				// clone failed; cleanup remnants
				this.os.RemoveAll(pkgPath)
				return nil, model.ErrPkgPullAuthorization{}
			case git.ErrRepositoryAlreadyExists.Error():
				return nil, nil
				// NoOp on repo already exists
			default:
				// clone failed; cleanup remnants
				this.os.RemoveAll(pkgPath)
				return nil, err
			}
		}

		// remove pkg '.git' sub dir
		return nil, this.os.RemoveAll(filepath.Join(pkgPath, ".git"))
	})

	return err

}
