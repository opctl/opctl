package pkg

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
)

// Pull pulls 'pkgRef' to 'path'
// returns ErrAuthenticationFailed on authentication failure
func (this _Pkg) Pull(
	path string,
	pkgRef *PkgRef,
	opts *PullOpts,
) error {

	cloneOptions := &git.CloneOptions{
		URL:           fmt.Sprintf("https://%v", pkgRef.FullyQualifiedName),
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%v", pkgRef.Version)),
		Depth:         1,
		Progress:      os.Stdout,
	}

	if nil != opts {
		cloneOptions.Auth = http.NewBasicAuth(opts.Username, opts.Password)
	}

	pkgPath := pkgRef.ToPath(path)

	if _, err := this.git.PlainClone(
		pkgPath,
		false,
		cloneOptions,
	); nil != err {
		switch err.Error() {
		// @TODO update to handle authentication & authorization errors separately once go-git does so
		case transport.ErrAuthorizationRequired.Error():
			// clone failed; cleanup remnants
			this.os.RemoveAll(pkgPath)
			return ErrAuthenticationFailed{}
		case git.ErrRepositoryAlreadyExists.Error():
		// NoOp on repo already exists
		default:
			// clone failed; cleanup remnants
			this.os.RemoveAll(pkgPath)
			return err
		}
	}

	return nil
}
