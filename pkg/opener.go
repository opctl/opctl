package pkg

//go:generate counterfeiter -o ./fakeOpener.go --fake-name fakeOpener ./ opener

import (
	"fmt"
)

// opener opens a package
type opener interface {
	Open(
		pkgRef string,
	) (
		handle,
		error,
	)
}

func newOpener(
	cachePath string,
) opener {
	return _opener{
		cachePath: cachePath,
		puller:    newPuller(),
		refParser: newRefParser(),
	}
}

type _opener struct {
	cachePath string
	puller    puller
	refParser refParser
	resolver  resolver
}

func (opener _opener) Open(
	pkgRef string,
) (
	handle,
	error,
) {
	parsedPkgRef, err := opener.refParser.ParseRef(pkgRef)
	if nil != err {
		return nil, err
	}

	var (
		pkgPath string
		ok      bool
	)
	pkgPath, ok = opener.resolver.Resolve(parsedPkgRef, opener.cachePath)
	if !ok {
		// pkg not resolved; attempt to pull it
		err := opener.puller.Pull(opener.cachePath, parsedPkgRef, nil)
		if nil != err {
			return nil, err
		}

		// resolve pulled pkg
		pkgPath, ok = opener.resolver.Resolve(parsedPkgRef, opener.cachePath)
		if !ok {
			return nil, fmt.Errorf("Unable to resolve pulled pkg '%v' from '%v'", pkgRef, opener.cachePath)
		}
	}

	return newHandleLocal(pkgPath), nil
}
