package pkg

//go:generate counterfeiter -o ./fakeRefResolver.go --fake-name fakeRefResolver ./ refResolver

import (
	"github.com/virtual-go/vos"
	pathPkg "path"
)

type refResolver interface {
	Resolve(pkgRef string) string
}

func newRefResolver(
	os vos.VOS,
) refResolver {
	return _refResolver{
		os: os,
	}
}

type _refResolver struct {
	os vos.VOS
}

func (this _refResolver) Resolve(
	pkgRef string,
) string {
	if localPkgRef, ok := this.resolveLocal(pkgRef); ok {
		return localPkgRef
	}
	// otherwise return unmodified
	return pkgRef
}

func (this _refResolver) resolveLocal(
	pkgRef string,
) (string, bool) {

	var (
		localPkgRef string
		wd          string
		err         error
	)
	if wd, err = this.os.Getwd(); nil != err {
		return "", false
	}

	dirName := pathPkg.Base(pathPkg.Dir(pkgRef))
	if dirName != DotOpspecDirName {
		// ensure resolved from .opspec dir
		localPkgRef = pathPkg.Join(DotOpspecDirName, pkgRef)
	}

	if _, err := this.os.Stat(localPkgRef); nil == err {
		// local ref
		if !pathPkg.IsAbs(localPkgRef) {
			// ensure path absolute
			localPkgRef = pathPkg.Join(wd, DotOpspecDirName, pkgRef)
		}
		return localPkgRef, true
	}
	return "", false
}
