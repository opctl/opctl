package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
)

// opInstall implements "op install" sub command
func opInstall(
	ctx context.Context,
	dataResolver dataresolver.DataResolver,
	opPath []string,
	opRef string,
	path string,
	creds *model.Creds,
) error {
	// install the whole pkg in case relative (intra pkg) refs exist
	opRefParts := strings.SplitN(opRef, "#", 2)
	var dataRef string
	if len(opRefParts) == 1 {
		dataRef = opRefParts[0]
	} else {
		if verAndPathParts := strings.SplitN(opRefParts[1], "/", 2); len(verAndPathParts) != 1 {
			dataRef = fmt.Sprintf("%s#%s", opRefParts[0], verAndPathParts[0])
		}
	}

	opDirHandle, err := dataResolver.Resolve(
		ctx,
		dataRef,
		opPath,
		creds,
	)
	if err != nil {
		return err
	}

	return opspec.Install(
		ctx,
		filepath.Join(path, dataRef),
		opDirHandle,
	)
}
