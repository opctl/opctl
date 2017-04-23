package pkg

//go:generate counterfeiter -o ./fakeGetter.go --fake-name fakeGetter ./ getter

import (
	"github.com/opspec-io/sdk-golang/model"
)

type getter interface {
	Get(
		basePath,
		pkgRef string,
	) (*model.PkgManifest, error)
}

func newGetter(
	manifestUnmarshaller manifestUnmarshaller,
	resolver resolver,
) getter {
	return _getter{
		manifestUnmarshaller: manifestUnmarshaller,
		resolver:             resolver,
	}
}

type _getter struct {
	manifestUnmarshaller manifestUnmarshaller
	resolver             resolver
}

func (this _getter) Get(
	basePath,
	pkgRef string,
) (*model.PkgManifest, error) {
	if localPkg, ok := this.resolver.Resolve(basePath, pkgRef); ok {
		return this.manifestUnmarshaller.Unmarshal(localPkg)
	}
	return nil, ErrPkgNotFound{}
}
