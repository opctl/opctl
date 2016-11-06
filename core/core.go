package core

//go:generate counterfeiter -o ./fakeCore.go --fake-name FakeCore ./ Core

import (
  "github.com/opspec-io/sdk-golang/pkg/models"
  "github.com/opspec-io/sdk-golang/pkg/bundle"
  "github.com/opspec-io/engine/pkg/containerengine"
  "github.com/opspec-io/engine/util/eventing"
  "github.com/opspec-io/engine/util/uniquestring"
  "github.com/opspec-io/engine/util/pathnormalizer"
)

type Core interface {
  GetEventStream(
  eventChannel chan models.Event,
  ) (err error)

  KillOpRun(
  req models.KillOpRunReq,
  ) (
  err error,
  )

  StartOpRun(
  req models.StartOpRunReq,
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

  opspecSdk := bundle.New()

  storage := newStorage()

  opRunner := newOpRunner(
    containerEngine,
    eventStream,
    opspecSdk,
    storage,
    uniqueStringFactory,
  )

  core = &_core{
    bundle:bundle.New(),
    containerEngine:containerEngine,
    eventStream:eventStream,
    opRunner:opRunner,
    pathNormalizer:pathnormalizer.NewPathNormalizer(),
    storage:storage,
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
  storage             storage
  uniqueStringFactory uniquestring.UniqueStringFactory
}
