package core

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Core

import (
	"github.com/opspec-io/opctl/util/containerprovider"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
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
	containerProvider containerprovider.ContainerProvider,
) (core Core) {
	uniqueStringFactory := uniquestring.NewUniqueStringFactory()

	pubSub := pubsub.New()

	_bundle := bundle.New()

	dcgNodeRepo := newDcgNodeRepo()

	caller := newCaller(
		newContainerCaller(
			_bundle,
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
		_bundle,
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
		pubSub:              pubSub,
		opCaller:            opCaller,
		dcgNodeRepo:         dcgNodeRepo,
		uniqueStringFactory: uniqueStringFactory,
	}

	return
}

type _core struct {
	containerProvider   containerprovider.ContainerProvider
	pubSub              pubsub.PubSub
	caller              caller
	dcgNodeRepo         dcgNodeRepo
	uniqueStringFactory uniquestring.UniqueStringFactory
	opCaller            opCaller
}
