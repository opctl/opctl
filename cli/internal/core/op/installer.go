package op

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
)

// Installer exposes the "op install" sub command
type Installer interface {
	Install(
		ctx context.Context,
		path,
		opRef,
		username,
		password string,
	) error
}

// newInstaller returns an initialized "op install" sub command
func newInstaller(
	dataResolver dataresolver.DataResolver,
) Installer {
	return _installer{
		dataResolver: dataResolver,
	}
}

type _installer struct {
	dataResolver dataresolver.DataResolver
}

func (ivkr _installer) Install(
	ctx context.Context,
	path,
	opRef,
	username,
	password string,
) error {
	// install the whole pkg in case relative (intra pkg) refs exist
	opRefParts := strings.SplitN(opRef, "#", 2)
	var pkgRef string
	if len(opRefParts) == 1 {
		pkgRef = opRefParts[0]
	} else {
		if verAndPathParts := strings.SplitN(opRefParts[1], "/", 2); len(verAndPathParts) != 1 {
			pkgRef = fmt.Sprintf("%s#%s", opRefParts[0], verAndPathParts[0])
		}
	}

	opDirHandle, err := ivkr.dataResolver.Resolve(
		pkgRef,
		&model.Creds{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		return err
	}

	return opspec.Install(
		ctx,
		filepath.Join(path, pkgRef),
		opDirHandle,
	)
}
