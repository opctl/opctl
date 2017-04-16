package pkg

//go:generate counterfeiter -o ./fakeGetter.go --fake-name fakeGetter ./ getter

import (
	"fmt"
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/vgit"
	"github.com/virtual-go/fs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
	"path"
	"strings"
)

type getter interface {
	Get(
		req *GetReq,
	) (*model.PkgManifest, error)
}

func newGetter(
	fs fs.FS,
	manifestUnmarshaller manifestUnmarshaller,
	refResolver refResolver,
) getter {
	return _getter{
		fs:                   fs,
		git:                  vgit.New(),
		manifestUnmarshaller: manifestUnmarshaller,
		refResolver:          refResolver,
	}
}

type _getter struct {
	fs                   fs.FS
	git                  vgit.VGit
	manifestUnmarshaller manifestUnmarshaller
	refResolver          refResolver
}

func (this _getter) Get(
	req *GetReq,
) (*model.PkgManifest, error) {
	resolvedPkgRef := this.refResolver.Resolve(req.PkgRef)
	if _, err := this.fs.Stat(resolvedPkgRef); nil == err {
		// observe default behavior; resolve from .opspec
		return this.manifestUnmarshaller.Unmarshal(resolvedPkgRef)
	}
	return this.getRemote(resolvedPkgRef, req.Username, req.Password)
}

func (this _getter) getRemote(
	pkgRef string,
	username string,
	password string,
) (*model.PkgManifest, error) {

	stringParts := strings.Split(pkgRef, "#")
	if len(stringParts) != 2 {
		return nil, fmt.Errorf(
			"Invalid remote pkgRef:'%v'. Valid remote pkgRef's are of the form: 'host/path#semver",
			pkgRef,
		)
	}
	repoName := stringParts[0]
	repoRefName := stringParts[1]

	gitPkgPath := path.Join(
		appdatapath.New().PerUser(),
		"opspec",
		"cache",
		"pkgs",
		repoName,
		repoRefName,
	)
	if _, err := this.fs.Stat(gitPkgPath); nil != err {
		// pkg not resolved on node; clone it
		cloneOptions := &git.CloneOptions{
			URL:           fmt.Sprintf("https://%v", repoName),
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%v", repoRefName)),
			Depth:         1,
			Progress:      os.Stdout,
		}

		if "" != username && "" != password {
			cloneOptions.Auth = http.NewBasicAuth(username, password)
		}

		err := this.git.PlainClone(gitPkgPath, false, cloneOptions)
		if nil != err {
			// clone failed; cleanup remnants
			this.fs.RemoveAll(gitPkgPath)
			return nil, err
		}
	}
	return this.manifestUnmarshaller.Unmarshal(gitPkgPath)
}
