package core

import (
	"github.com/opctl/opctl/util/containerprovider"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opctl/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/containercall"
	opspecNodeCorePkg "github.com/opspec-io/sdk-golang/node/core"
)

func New(
	pubSub pubsub.PubSub,
	containerProvider containerprovider.ContainerProvider,
	rootFSPath string,
) (core opspecNodeCorePkg.Core) {
	uniqueStringFactory := uniquestring.NewUniqueStringFactory()

	dcgNodeRepo := newDCGNodeRepo()

	opKiller := newOpKiller(dcgNodeRepo, containerProvider)

	caller := newCaller(
		newContainerCaller(
			containerProvider,
			containercall.New(rootFSPath),
			pubSub,
			dcgNodeRepo,
		),
	)

	caller.setParallelCaller(
		newParallelCaller(
			caller,
			opKiller,
			pubSub,
			uniqueStringFactory,
		),
	)

	caller.setSerialCaller(
		newSerialCaller(
			caller,
			pubSub,
			uniqueStringFactory,
		),
	)

	opCaller := newOpCaller(
		pubSub,
		dcgNodeRepo,
		caller,
		rootFSPath,
	)

	caller.setOpCaller(
		opCaller,
	)

	core = _core{
		containerProvider:   containerProvider,
		dcgNodeRepo:         dcgNodeRepo,
		opCaller:            opCaller,
		opKiller:            opKiller,
		pubSub:              pubSub,
		uniqueStringFactory: uniqueStringFactory,
	}

	return
}

type _core struct {
	containerProvider   containerprovider.ContainerProvider
	dcgNodeRepo         dcgNodeRepo
	opCaller            opCaller
	opKiller            opKiller
	pubSub              pubsub.PubSub
	uniqueStringFactory uniquestring.UniqueStringFactory
}
