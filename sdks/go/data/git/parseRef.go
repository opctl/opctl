package git

import (
	"errors"
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
	if err != nil {
		return nil, err
	}

	// fragment MAY be in format: SEM_VER/OP_PATH
	version := strings.SplitN(refURI.Fragment, "/", 2)[0]
	if version == "" {
		return nil, errors.New("missing version")
	}

	return &ref{
		Name:    path.Join(refURI.Host, refURI.Path),
		Version: version,
	}, nil
}
