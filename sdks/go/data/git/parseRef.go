package git

import (
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

// parseRef string to object
func parseRef(
	dataRef string,
) (*ref, error) {
	refURI, err := url.Parse(filepath.ToSlash(dataRef))
	if nil != err {
		return nil, err
	}

	return &ref{
		Name: path.Join(refURI.Host, refURI.Path),
		// fragment MAY be in format: SEM_VER/OP_PATH
		Version: strings.SplitN(refURI.Fragment, "/", 2)[0],
	}, nil
}
