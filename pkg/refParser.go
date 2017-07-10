package pkg

//go:generate counterfeiter -o ./fakeRefParser.go --fake-name fakeRefParser ./ refParser

import (
	"net/url"
	"path"
	"path/filepath"
)

// refParser parses pkg refs
type refParser interface {
	Parse(
		pkgRef string,
	) (
		*Ref,
		error,
	)
}

func newRefParser() refParser {
	return _refParser{}
}

type _refParser struct{}

type Ref struct {
	Name    string
	Version string
}

// ToPath constructs a filesystem path for a Ref, assuming the provided base path
func (pr Ref) ToPath(basePath string) string {
	return filepath.Join(basePath, filepath.FromSlash(pr.Name), pr.Version)
}

// Parse parses a pkgRef
func (rp _refParser) Parse(
	pkgRef string,
) (*Ref, error) {
	refURI, err := url.Parse(filepath.ToSlash(pkgRef))
	if nil != err {
		return nil, err
	}

	return &Ref{
		Name:    path.Join(refURI.Host, refURI.Path),
		Version: refURI.Fragment,
	}, nil
}
