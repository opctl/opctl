package git

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
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

	if version != "" {
		if _, err = semver.Parse(version); err != nil {
			err = fmt.Errorf("%s is not a valid semver", version)
		}
	}

	return &ref{
		Name:    path.Join(refURI.Host, refURI.Path),
		Version: version,
	}, err
}
