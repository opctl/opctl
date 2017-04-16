package pkg

//go:generate counterfeiter -o ./fakeGetter.go --fake-name fakeGetter ./ getter

import (
	"errors"
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
) getter {
	return _getter{
		fs:                   fs,
		git:                  vgit.New(),
		manifestUnmarshaller: manifestUnmarshaller,
	}
}

type _getter struct {
	fs                   fs.FS
	git                  vgit.VGit
	manifestUnmarshaller manifestUnmarshaller
}

func (this _getter) Get(
	req *GetReq,
) (*model.PkgManifest, error) {

	if _, err := this.fs.Stat(req.PkgRef); nil == err {
		return this.manifestUnmarshaller.Unmarshal(req.PkgRef)
	}
	return this.getRemote(req)
}

func (this _getter) getRemote(
	req *GetReq,
) (*model.PkgManifest, error) {

	stringParts := strings.Split(req.PkgRef, "#")
	if len(stringParts) != 2 {
		return nil, errors.New("Invalid pkgRef, version not provided")
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

		if "" != req.Username && "" != req.Password {
			cloneOptions.Auth = http.NewBasicAuth(req.Username, req.Password)
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
