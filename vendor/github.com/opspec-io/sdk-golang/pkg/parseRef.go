package pkg

import (
	"net/url"
	"path/filepath"
)

type PkgRef struct {
	FullyQualifiedName string
	Version            string
}

// ParseRef parses a pkgRef
func (p _Pkg) ParseRef(
	pkgRef string,
) (*PkgRef, error) {
	pkgRefURI, err := url.Parse(pkgRef)
	if nil != err {
		return nil, err
	}

	return &PkgRef{
		FullyQualifiedName: filepath.Join(pkgRefURI.Host, pkgRefURI.Path),
		Version:            pkgRefURI.Fragment,
	}, nil
}
