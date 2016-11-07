package core

//go:generate counterfeiter -o ./fakeCore.go --fake-name FakeCore ./ Core

import (
  "github.com/opspec-io/sdk-golang/pkg/model"
  "github.com/opspec-io/sdk-golang/pkg/bundle"
  "github.com/opspec-io/engine/pkg/containerengine"
  "github.com/opspec-io/engine/util/eventing"
  "github.com/opspec-io/engine/util/uniquestring"
  "github.com/opspec-io/engine/util/pathnormalizer"
)

type Core interface {
  GetEventStream(
  eventChannel chan model.Event,
  ) (err error)

  KillOpRun(
  req model.KillOpRunReq,
  )

  StartOpRun(
  req model.StartOpRunReq,
  ) (
  opRunId string,
  err error,
  )
}

func New(
containerEngine containerengine.ContainerEngine,
) (core Core) {

  /* factories */
  uniqueStringFactory := uniquestring.NewUniqueStringFactory()

  /* components */
  eventStream := eventing.NewEventStream()

  _bundle := bundle.New()

  opRunRepo := newOpRunRepo()

  opRunner := newOpRunner(
    containerEngine,
    eventStream,
    _bundle,
    opRunRepo,
    uniqueStringFactory,
  )

  core = &_core{
    bundle:_bundle,
    containerEngine:containerEngine,
    eventStream:eventStream,
    opRunner:opRunner,
    pathNormalizer:pathnormalizer.NewPathNormalizer(),
    opRunRepo:opRunRepo,
    uniqueStringFactory:uniqueStringFactory,
  }

  return
}

type _core struct {
  bundle              bundle.Bundle
  containerEngine     containerengine.ContainerEngine
  eventStream         eventing.EventStream
  opRunner            opRunner
  pathNormalizer      pathnormalizer.PathNormalizer
  opRunRepo           opRunRepo
  uniqueStringFactory uniquestring.UniqueStringFactory
}
