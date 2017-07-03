package pkg

import (
	"net/url"
	"path"
	"path/filepath"
)

// refParser parses pkg refs
type refParser interface {
	ParseRef(
		pkgRef string,
	) (
		*PkgRef,
		error,
	)
}

func newRefParser() refParser {
	return _refParser{}
}

type _refParser struct{}

type PkgRef struct {
	FullyQualifiedName string
	Version            string
}

// ToPath constructs a filesystem path for a PkgRef, assuming the provided base path
func (pr PkgRef) ToPath(basePath string) string {
	return filepath.Join(basePath, filepath.FromSlash(pr.FullyQualifiedName), pr.Version)
}

// ParseRef parses a pkgRef
func (rp _refParser) ParseRef(
	pkgRef string,
) (*PkgRef, error) {
	refURI, err := url.Parse(filepath.ToSlash(pkgRef))
	if nil != err {
		return nil, err
	}

	return &PkgRef{
		FullyQualifiedName: path.Join(refURI.Host, refURI.Path),
		Version:            refURI.Fragment,
	}, nil
}
