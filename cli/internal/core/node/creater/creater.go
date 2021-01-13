package creater

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"

	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/node/core"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/k8s"
)

// Creater exposes the "node create" sub command
//counterfeiter:generate -o fakes/creater.go . Creater
type Creater interface {
	Create(
		opts local.NodeCreateOpts,
	) error
}

// New returns an initialized "node create" command
func New() Creater {
	return _creater{}
}

type _creater struct{}

func (ivkr _creater) Create(
	opts local.NodeCreateOpts,
) error {
	dataDir, err := datadir.New(opts.DataDir)
	if nil != err {
		return err
	}

	if err := dataDir.InitAndLock(); nil != err {
		return err
	}

	var containerRuntime containerruntime.ContainerRuntime
	if "k8s" == opts.ContainerRuntime {
		containerRuntime, err = k8s.New()
		if nil != err {
			return err
		}
	} else {
		containerRuntime, err = docker.New()
		if nil != err {
			return err
		}
	}

	err = newHTTPListener(
		core.New(
			containerRuntime,
			dataDir.Path(),
		),
	).
		Listen(
			context.Background(),
			opts.ListenAddress,
		)

	if nil != err {
		return err
	}

	return nil
}
