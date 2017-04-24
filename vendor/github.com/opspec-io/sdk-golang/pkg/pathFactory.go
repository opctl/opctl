package pkg

import (
	"fmt"
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	"path"
	"strings"
)

// constructCachePath constructs the cache path of a pkg from a pkgRef
func constructCachePath(pkgRef string) (string, error) {
	stringParts := strings.Split(pkgRef, "#")
	if len(stringParts) != 2 {
		return "", fmt.Errorf(
			"Invalid remote pkgRef: '%v'. Valid remote pkgRef's are of the form: 'host/path#semver",
			pkgRef,
		)
	}
	repoName := stringParts[0]
	repoRefName := stringParts[1]

	return path.Join(
		appdatapath.New().PerUser(),
		"opspec",
		"cache",
		"pkgs",
		repoName,
		repoRefName,
	), nil
}
