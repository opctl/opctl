package core

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Core

import (
	"github.com/opctl/opctl/util/containerprovider"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opctl/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/validate"
)

type Core interface {
	GetEventStream(
		req *model.GetEventStreamReq,
		eventChannel chan *model.Event,
	) error

	KillOp(
		req model.KillOpReq,
	)

	StartOp(
		req model.StartOpReq,
	) (
		callId string,
		err error,
	)
}

func New(
	pubSub pubsub.PubSub,
	containerProvider containerprovider.ContainerProvider,
) (core Core) {
	uniqueStringFactory := uniquestring.NewUniqueStringFactory()

	dcgNodeRepo := newDCGNodeRepo()

	opKiller := newOpKiller(dcgNodeRepo, containerProvider)

	caller := newCaller(
		newContainerCaller(
			containerProvider,
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
		pkg.New(),
		pubSub,
		dcgNodeRepo,
		caller,
		uniqueStringFactory,
		validate.New(),
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
