package pkg

import (
	"fmt"
	"github.com/golang-interfaces/gopkg.in-src-d-go-git.v4"
	"github.com/golang-interfaces/ios"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
	"path/filepath"
)

type puller interface {
	// Pull pulls 'pkgRef' to 'path'
	// nil pullCreds will be ignored
	// returns ErrAuthenticationFailed on authentication failure
	Pull(
		path string,
		pkgRef string,
		pullCreds *PullCreds,
	) error
}

func newPuller() puller {
	return _puller{
		git: igit.New(),
		os:  ios.New(),
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
	authOpts *PullCreds,
) error {

	parsedPkgRef, err := this.refParser.Parse(pkgRef)
	if nil != err {
		return err
	}

	cloneOptions := &git.CloneOptions{
		URL:           fmt.Sprintf("https://%v", parsedPkgRef.Name),
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%v", parsedPkgRef.Version)),
		Depth:         1,
		Progress:      os.Stdout,
	}

	if nil != authOpts {
		cloneOptions.Auth = http.NewBasicAuth(authOpts.Username, authOpts.Password)
	}

	pkgPath := parsedPkgRef.ToPath(path)

	if _, err := this.git.PlainClone(
		pkgPath,
		false,
		cloneOptions,
	); nil != err {
		switch err.Error() {
		// @TODO update to Handle authentication & authorization errors separately once go-git does so
		case transport.ErrAuthorizationRequired.Error():
			// clone failed; cleanup remnants
			this.os.RemoveAll(pkgPath)
			return ErrAuthenticationFailed{}
		case git.ErrRepositoryAlreadyExists.Error():
			return nil
		// NoOp on repo already exists
		default:
			// clone failed; cleanup remnants
			this.os.RemoveAll(pkgPath)
			return err
		}
	}

	// remove pkg '.git' sub dir
	return this.os.RemoveAll(filepath.Join(pkgPath, ".git"))
}
