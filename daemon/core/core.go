package core

//go:generate counterfeiter -o ./fakeCore.go --fake-name FakeCore ./ Core

import (
	"github.com/opspec-io/opctl/pkg/containerengine"
	"github.com/opspec-io/opctl/util/eventbus"
	"github.com/opspec-io/opctl/util/pathnormalizer"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/pkg/validate"
)

type Core interface {
	GetEventStream(
		req *model.GetEventStreamReq,
		eventChannel chan model.Event,
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
	containerEngine containerengine.ContainerEngine,
) (core Core) {
	uniqueStringFactory := uniquestring.NewUniqueStringFactory()

	eventBus := eventbus.New()

	_bundle := bundle.New()

	nodeRepo := newNodeRepo()

	caller := newCaller(
		newContainerCaller(
			_bundle,
			containerEngine,
			eventBus,
			nodeRepo,
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
		eventBus,
		nodeRepo,
		caller,
		uniqueStringFactory,
		validate.New(),
	)

	caller.setOpCaller(
		opCaller,
	)

	core = _core{
		containerEngine:     containerEngine,
		eventBus:            eventBus,
		opCaller:            opCaller,
		pathNormalizer:      pathnormalizer.NewPathNormalizer(),
		nodeRepo:            nodeRepo,
		uniqueStringFactory: uniqueStringFactory,
	}

	return
}

type _core struct {
	containerEngine     containerengine.ContainerEngine
	eventBus            eventbus.EventBus
	caller              caller
	pathNormalizer      pathnormalizer.PathNormalizer
	nodeRepo            nodeRepo
	uniqueStringFactory uniquestring.UniqueStringFactory
	opCaller            opCaller
}
