package creater

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"

	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/model"
	"github.com/opctl/opctl/sdks/go/node/core"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

// Creater exposes the "node create" sub command
//counterfeiter:generate -o fakes/creater.go . Creater
type Creater interface {
	Create(
		opts model.NodeCreateOpts,
	)
}

// New returns an initialized "node create" command
func New() Creater {
	return _creater{}
}

type _creater struct{}

func (ivkr _creater) Create(
	opts model.NodeCreateOpts,
) {
	dataDir, err := datadir.New(opts.DataDir)
	if nil != err {
		panic(err)
	}

	if err := dataDir.InitAndLock(); nil != err {
		panic(err)
	}

	containerRuntime, err := docker.New()
	if nil != err {
		panic(err)
	}

	err = newHTTPListener(
		core.New(
			pubsub.New(
				pubsub.NewBadgerDBEventStore(
					dataDir.EventDBPath(),
				),
			),
			containerRuntime,
			dataDir.Path(),
		),
	).
		Listen(
			context.Background(),
			opts.ListenAddress,
		)

	if nil != err {
		panic(err)
	}

}
