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
	localResolver localResolver,
) getter {
	return _getter{
		fs:                   fs,
		git:                  vgit.New(),
		manifestUnmarshaller: manifestUnmarshaller,
		localResolver:        localResolver,
	}
}

type _getter struct {
	fs                   fs.FS
	git                  vgit.VGit
	manifestUnmarshaller manifestUnmarshaller
	localResolver        localResolver
}

func (this _getter) Get(
	req *GetReq,
) (*model.PkgManifest, error) {
	if localPkg, ok := this.localResolver.Resolve(req.BasePath, req.PkgRef); ok {
		return this.manifestUnmarshaller.Unmarshal(localPkg)
	}
	return this.getRemote(req.PkgRef, req.Username, req.Password)
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
	return this.manifestUnmarshaller.Unmarshal(gitPkgPath)
}
