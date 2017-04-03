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
	pathPkg "path"
	"strings"
)

type getter interface {
	Get(
		req *GetReq,
	) (packageView *model.PackageView, err error)
}

func newGetter(
	fs fs.FS,
	viewFactory viewFactory,
) getter {
	return _getter{
		fs:          fs,
		git:         vgit.New(),
		viewFactory: viewFactory,
	}
}

type _getter struct {
	fs          fs.FS
	git         vgit.VGit
	viewFactory viewFactory
}

func (this _getter) Get(
	req *GetReq,
) (packageView *model.PackageView, err error) {
	embeddedPkgPath := pathPkg.Join(req.path, ".opspec", req.pkgRef)
	if _, err = this.fs.Stat(embeddedPkgPath); nil == err {
		return this.viewFactory.Construct(embeddedPkgPath)
	}
	return this.getRemote(req)
}

func (this _getter) getRemote(
	req *GetReq,
) (packageView *model.PackageView, err error) {

	stringParts := strings.Split(req.pkgRef, "#")
	repoName := stringParts[0]
	repoRefName := stringParts[1]

	gitPkgPath := pathPkg.Join(
		appdatapath.New().Global(),
		".opspec",
		"refs",
		repoName,
		repoRefName,
	)

	if _, err = this.fs.Stat(gitPkgPath); nil != err {
		// ref not resolved on node; pull it
		cloneOptions := &git.CloneOptions{
			URL:           fmt.Sprintf("https://%v", repoName),
			ReferenceName: plumbing.ReferenceName(repoRefName),
			SingleBranch:  true,
		}

		if "" != req.username && "" != req.password {
			cloneOptions.Auth = http.NewBasicAuth(req.username, req.password)
		}

		_, err = git.PlainClone("", false, cloneOptions)
		if nil != err {
			return
		}
	}
	packageView, err = this.viewFactory.Construct(gitPkgPath)

	return
}
