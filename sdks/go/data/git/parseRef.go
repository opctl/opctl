package git

import (
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
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
		return nil, model.ErrDataMissingVersion{}
	}

	return &ref{
		Name:    path.Join(refURI.Host, refURI.Path),
		Version: version,
	}, nil
}
