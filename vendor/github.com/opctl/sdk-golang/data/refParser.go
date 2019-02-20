package data

//go:generate counterfeiter -o ./fakeRefParser.go --fake-name fakeRefParser ./ refParser

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

// refParser parses "dataRef"
type refParser interface {
	Parse(
		dataRef string,
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
	OpPath  string
}

// ToPath constructs a filesystem path for a Ref, assuming the provided base path
func (pr Ref) ToPath(basePath string) string {
	crossPlatPath := filepath.FromSlash(fmt.Sprintf("%v#%v", pr.Name, pr.Version))
	return filepath.Join(basePath, crossPlatPath)
}

// Parse parses a ref
func (rp _refParser) Parse(
	dataRef string,
) (*Ref, error) {
	refURI, err := url.Parse(filepath.ToSlash(dataRef))
	if nil != err {
		return nil, err
	}

	return &Ref{
		Name: path.Join(refURI.Host, refURI.Path),
		// fragment MAY be in format: SEM_VER/OP_PATH
		Version: strings.SplitN(refURI.Fragment, "/", 2)[0],
	}, nil
}
