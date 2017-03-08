package core

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Core

import (
	"github.com/opspec-io/opctl/util/containerprovider"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/managepackages"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/pkg/validate"
)

type Core interface {
	GetEventStream(
		req *model.GetEventStreamReq,
		eventChannel chan *model.Event,
	) (err error)

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

	dcgNodeRepo := newDcgNodeRepo()

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
			uniqueStringFactory,
		),
	)

	caller.setSerialCaller(
		newSerialCaller(
			caller,
			uniqueStringFactory,
		),
	)

	opCaller := newOpCaller(
		managepackages.New(),
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
		pubSub:              pubSub,
		uniqueStringFactory: uniqueStringFactory,
	}

	return
}

type _core struct {
	containerProvider   containerprovider.ContainerProvider
	dcgNodeRepo         dcgNodeRepo
	opCaller            opCaller
	pubSub              pubsub.PubSub
	uniqueStringFactory uniquestring.UniqueStringFactory
}
