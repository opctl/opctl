package core

//go:generate counterfeiter -o ./fakeCore.go --fake-name FakeCore ./ Core

import (
	"github.com/opspec-io/engine/pkg/containerengine"
	"github.com/opspec-io/engine/util/eventbus"
	"github.com/opspec-io/engine/util/pathnormalizer"
	"github.com/opspec-io/engine/util/uniquestring"
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

	/* factories */
	uniqueStringFactory := uniquestring.NewUniqueStringFactory()

	/* components */
	eventBus := eventbus.NewEventBus()

	_bundle := bundle.New()

	nodeRepo := newNodeRepo()

	opOrchestrator := &_opOrchestrator{
		bundle:              _bundle,
		eventBus:            eventBus,
		nodeRepo:            nodeRepo,
		uniqueStringFactory: uniqueStringFactory,
		validate:            validate.New(),
	}

	opOrchestrator.orchestrator = newOrchestrator(
		_bundle,
		containerEngine,
		eventBus,
		nodeRepo,
		opOrchestrator,
		uniqueStringFactory,
	)

	core = &_core{
		containerEngine:     containerEngine,
		eventBus:            eventBus,
		opOrchestrator:      opOrchestrator,
		pathNormalizer:      pathnormalizer.NewPathNormalizer(),
		nodeRepo:            nodeRepo,
		uniqueStringFactory: uniqueStringFactory,
	}

	return
}

type _core struct {
	containerEngine     containerengine.ContainerEngine
	eventBus            eventbus.EventBus
	orchestrator        orchestrator
	pathNormalizer      pathnormalizer.PathNormalizer
	nodeRepo            nodeRepo
	uniqueStringFactory uniquestring.UniqueStringFactory
	opOrchestrator      opOrchestrator
}
