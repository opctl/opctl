package main

import (
	"context"
	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/node/core"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/k8s"
)

// node command
func node(
  ctx context.Context,
	nodeCreateOpts local.NodeCreateOpts,
) error {
	dataDir, err := datadir.New(nodeCreateOpts.DataDir)
	if nil != err {
		return err
	}

	if err := dataDir.InitAndLock(); nil != err {
		return err
	}

	var containerRuntime containerruntime.ContainerRuntime
	if "k8s" == nodeCreateOpts.ContainerRuntime {
		containerRuntime, err = k8s.New()
		if nil != err {
			return err
		}
	} else {
		containerRuntime, err = docker.New(ctx)
		if nil != err {
			return err
		}
	}

	return newHTTPListener(
		core.New(
      ctx,
			containerRuntime,
			dataDir.Path(),
		),
	).
		listen(
			ctx,
			nodeCreateOpts.ListenAddress,
		)

}
