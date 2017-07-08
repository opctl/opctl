package pkg

import (
	"github.com/golang-interfaces/ios"
	"os"
	"path/filepath"
)

type Resolver interface {
	// Resolve attempts to resolve a package according to opspec package resolution rules
	// nil opts will be ignored
	// returns ErrAuthenticationFailed on authentication failure
	Resolve(
		pkgRef string,
		opts *ResolveOpts,
	) (
		Handle,
		error,
	)
}

func newResolver(
	cachePath string,
) Resolver {
	return _resolver{
		cachePath: cachePath,
		os:        ios.New(),
		puller:    newPuller(),
	}
}

type _resolver struct {
	cachePath string
	os        ios.IOS
	puller    puller
}

func (this _resolver) Resolve(
	pkgRef string,
	opts *ResolveOpts,
) (
	Handle,
	error,
) {
	var (
		lookPaths []string = []string{this.cachePath}
		pullCreds *PullCreds
	)
	if nil != opts {
		if "" != opts.BasePath {
			// resolve from provided base path first
			lookPaths = []string{opts.BasePath, this.cachePath}
		}
		pullCreds = opts.PullCreds
	}

	// attempt local resolve
	var resolvedPkgRef string
	if filepath.IsAbs(pkgRef) {
		if _, err := this.os.Stat(pkgRef); nil != err && !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}
		resolvedPkgRef = pkgRef
	} else {
		for _, lookPath := range lookPaths {

			// attempt to resolve from lookPath/.opspec dir
			testPath := filepath.Join(lookPath, DotOpspecDirName, pkgRef)
			_, err := this.os.Stat(testPath)
			if err == nil {
				// record & break on resolve
				resolvedPkgRef = testPath
				break
			}
			if nil != err && !os.IsNotExist(err) {
				// return actual errors
				return nil, err
			}

			// attempt to resolve from lookPath
			testPath = filepath.Join(lookPath, pkgRef)
			_, err = this.os.Stat(testPath)
			if err == nil {
				// record & break on resolve
				resolvedPkgRef = testPath
				break
			}
			if nil != err && !os.IsNotExist(err) {
				// return actual errors
				return nil, err
			}
		}
	}

	// attempt pull if local resolve failed
	if "" == resolvedPkgRef {
		err := this.puller.Pull(this.cachePath, pkgRef, pullCreds)
		if nil != err {
			return nil, err
		}
		resolvedPkgRef = filepath.Join(this.cachePath, pkgRef)
	}

	return newLocalHandle(resolvedPkgRef), nil
}
