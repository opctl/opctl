package pkg

import (
	"fmt"
	"strings"
)

type PkgRef struct {
	FullyQualifiedName string
	Version            string
}

func parsePkgRef(
	pkgRef string,
) (*PkgRef, error) {
	stringParts := strings.Split(pkgRef, "#")
	if len(stringParts) != 2 {
		return nil, fmt.Errorf(
			"Invalid remote pkgRef: '%v'. Valid remote pkgRef's are of the form: 'host/path#semver",
			pkgRef,
		)
	}

	return &PkgRef{
		FullyQualifiedName: stringParts[0],
		Version:            stringParts[1],
	}, nil
}
