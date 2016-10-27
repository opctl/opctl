package core

//go:generate counterfeiter -o ./fakeApi.go --fake-name FakeApi ./ Api

import (
  "github.com/opspec-io/sdk-golang/models"
  "github.com/opspec-io/sdk-golang/pkg/bundle"
)

type Api interface {
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
containerEngine ContainerEngine,
) (api Api) {

  /* factories */
  uniqueStringFactory := newUniqueStringFactory()

  /* components */
  eventStream := newEventStream()

  opspecSdk :=  bundle.New()

  storage := newStorage()

  opRunner := newOpRunner(
    containerEngine,
    eventStream,
    eventStream,
    opspecSdk,
    storage,
    uniqueStringFactory,
  )

  api = &_api{
    bundle:bundle.New(),
    containerEngine:containerEngine,
    eventStream:eventStream,
    opRunner:opRunner,
    pathNormalizer:newPathNormalizer(),
    storage:storage,
    uniqueStringFactory:uniqueStringFactory,
  }

  return
}

type _api struct {
  bundle bundle.Bundle
  containerEngine ContainerEngine
  eventStream eventStream
  opRunner opRunner
  pathNormalizer pathNormalizer
  storage storage
  uniqueStringFactory uniqueStringFactory
}
