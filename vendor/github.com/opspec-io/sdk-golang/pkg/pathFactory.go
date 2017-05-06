package pkg

import (
	"fmt"
	"github.com/appdataspec/sdk-golang/appdatapath"
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

	perUserAppDataPath, err := appdatapath.New().PerUser()
	if nil != err {
		return "", err
	}

	return path.Join(
		perUserAppDataPath,
		"opspec",
		"cache",
		"pkgs",
		repoName,
		repoRefName,
	), nil
}
