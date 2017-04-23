package pkg

import (
	"fmt"
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
	"path"
	"strings"
)

// Pull pulls a package from a remote source
func (this pkg) Pull(
	pkgRef string,
	opts *PullOpts,
) error {
	stringParts := strings.Split(pkgRef, "#")
	if len(stringParts) != 2 {
		return fmt.Errorf(
			"Invalid remote pkgRef: '%v'. Valid remote pkgRef's are of the form: 'host/path#semver",
			pkgRef,
		)
	}
	repoName := stringParts[0]
	repoRefName := stringParts[1]

	pkgPath := path.Join(
		appdatapath.New().PerUser(),
		"opspec",
		"cache",
		"pkgs",
		repoName,
		repoRefName,
	)
	cloneOptions := &git.CloneOptions{
		URL:           fmt.Sprintf("https://%v", repoName),
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%v", repoRefName)),
		Depth:         1,
		Progress:      os.Stdout,
	}

	if nil != opts {
		cloneOptions.Auth = http.NewBasicAuth(opts.Username, opts.Password)
	}

	err := this.git.PlainClone(pkgPath, false, cloneOptions)
	if nil != err {
		switch err.Error() {
		case git.ErrRepositoryAlreadyExists.Error():
		// NoOp on repo already exists
		default:
			// clone failed; cleanup remnants
			this.fs.RemoveAll(pkgPath)
			return err
		}
	}
	return nil
}
