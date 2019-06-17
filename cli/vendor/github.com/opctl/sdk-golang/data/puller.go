package data

//go:generate counterfeiter -o ./fakePuller.go --fake-name fakePuller ./ puller

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-interfaces/gopkg.in-src-d-go-git.v4"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/model"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

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

	if _, err := plr.git.PlainClone(
		opPath,
		false,
		cloneOptions,
	); nil != err {
		switch err.Error() {
		case transport.ErrAuthenticationRequired.Error():
			// clone failed; cleanup remnants
			plr.os.RemoveAll(opPath)
			return model.ErrDataProviderAuthentication{}
		case transport.ErrAuthorizationFailed.Error():
			// clone failed; cleanup remnants
			plr.os.RemoveAll(opPath)
			return model.ErrDataProviderAuthorization{}
		case git.ErrRepositoryAlreadyExists.Error():
			return nil
			// NoOp on repo already exists
		default:
			// clone failed; cleanup remnants
			plr.os.RemoveAll(opPath)
			return err
		}
	}

	// remove pkg '.git' sub dir
	return plr.os.RemoveAll(filepath.Join(opPath, ".git"))

}
