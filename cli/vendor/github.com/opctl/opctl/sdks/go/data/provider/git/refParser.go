package git

import (
	"github.com/opctl/opctl/sdks/go/data/provider/git/internal"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

// refParser parses "dataRef"
//counterfeiter:generate -o internal/fakes/refParser.go . refParser
type refParser interface {
	Parse(
		dataRef string,
	) (
		*internal.Ref,
		error,
	)
}

func newRefParser() refParser {
	return _refParser{}
}

type _refParser struct{}

// Parse parses a ref
func (rp _refParser) Parse(
	dataRef string,
) (*internal.Ref, error) {
	refURI, err := url.Parse(filepath.ToSlash(dataRef))
	if nil != err {
		return nil, err
	}

	return &internal.Ref{
		Name: path.Join(refURI.Host, refURI.Path),
		// fragment MAY be in format: SEM_VER/OP_PATH
		Version: strings.SplitN(refURI.Fragment, "/", 2)[0],
	}, nil
}
